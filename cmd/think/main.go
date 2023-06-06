// Package main implements a command line utility called "think", which provides an interactive shell interface
// for the user to describe tasks they want to execute in the shell. The tool leverages a language model (like GPT-4)
// to interpret these task descriptions and generate corresponding bash code. This code can then be edited by the user
// and subsequently executed in the shell.
//
// The utility supports two modes of operation:
//
//   - In the 'thinking' state, the user interacts with the AI model, describing the task they want to execute. The AI
//     responds with a concise explanation and a proposed bash command that corresponds to the described task.
//
//   - In the 'executing' state, the user reviews the proposed bash code, potentially editing it for correctness or
//     to better suit their specific needs.
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ryszard/agency/agent"
	"github.com/ryszard/agency/client"
	"github.com/ryszard/agency/client/openai"
	"github.com/sirupsen/logrus"
)

var (
	model      = flag.String("model", "", "model to use. The default is gpt-4")
	sendOutput = flag.Bool("send-output", false, "send the command you run and part of its stdout and stderr to the AI model. Useful, but potentially dangerous")
	logLevel   = flag.String("log-level", "error", "log level. One of: debug, info, warn, error, fatal, panic")

	flagsSet = make(map[string]bool)
)

func init() {
	flag.Usage = usage

	flag.Parse()

	// Set flagsSet for each flag that was provided in command line
	flag.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})

	// Set default values only if flag wasn't provided
	if !flagsSet["model"] {
		// Check for THINK_MODEL environment variable if the flag wasn't provided
		if envModel, ok := os.LookupEnv("THINK_MODEL"); ok {
			*model = envModel
		} else {
			*model = "gpt-4" // default value
		}
	}

	// Check for THINK_SEND_OUTPUT environment variable if the flag wasn't provided
	if !flagsSet["send-output"] {
		if envSendOutput, ok := os.LookupEnv("THINK_SEND_OUTPUT"); ok {
			// Parsing the string to bool. In case of error, it will be false.
			parsedEnvSendOutput, err := strconv.ParseBool(envSendOutput)
			if err != nil {
				log.Fatalf("error parsing THINK_SEND_OUTPUT environment variable: %v", err)
			}
			*sendOutput = parsedEnvSendOutput
		}
	}
}

//go:embed system.md
var SystemPrompt string

//go:embed user.md
var UserPrompt string

func usage() {
	fmt.Fprintf(os.Stderr, `think is a command-line tool that uses AI to generate and execute bash commands.

	Usage:
	
	  think [-model model] [-send-output] "your command"
	
	Options:
	
	  -model        Specifies the AI model to use. The default model is 'gpt-4'. 
	
	  -send-output  Send the command you run and part of its stdout and stderr to the AI model. 
					Useful, but potentially dangerous.
	
	You can also set the AI model to use with the 'THINK_MODEL' environment variable. If neither is provided, 
	the default model is 'gpt-4'. 
	
	Examples:
	
	  think "list all files in this directory"
	  think -model=gpt-4 "create a new directory called test"
	  think -send-output "print the contents of this file"
	`)
	flag.PrintDefaults()
}

func main() {

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetLevel(level)

	var cl client.Client = openai.New(os.Getenv("OPENAI_API_KEY"))
	cl = client.Retrying(cl, 1*time.Second, 5*time.Second, 10)

	if *model == "" {
		m, ok := os.LookupEnv("THINK_MODEL")
		if !ok {
			*model = "gpt-4"
		} else {
			*model = m
		}
	}

	ag := agent.New("scripter",
		agent.WithClient(cl),
		agent.WithModel(*model),
		agent.WithMaxTokens(500),
		agent.WithMemory(agent.TokenBufferMemory(3000)),
	)
	ag, err = agent.Templated(ag, map[string]string{
		"system": SystemPrompt,
		"user":   UserPrompt,
	})
	if err != nil {
		log.Fatal(err)
	}

	shellPath, ok := os.LookupEnv("SHELL")
	if !ok {
		shellPath = "/bin/bash"
	}

	operatingSystem := runtime.GOOS

	_, err = ag.System("system", struct {
		Shell string
		OS    string
	}{Shell: shellPath, OS: operatingSystem})
	if err != nil {
		log.Fatal(err)
	}

	repl, err := NewREPL(ag, shellPath, strings.Join(flag.Args(), " "), *sendOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer repl.Close()

	repl.Run()

}
