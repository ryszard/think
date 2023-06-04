// Package main implements a command line utility called "think", which provides an interactive shell interface
// for the user to describe tasks they want to execute in the shell. The tool leverages a language model (like GPT-4)
// to interpret these task descriptions and generate corresponding bash code. This code can then be edited by the user
// and subsequently executed in the shell.
//
// The utility supports multiple modes of operation, which are controlled by its current state:
//
//   - In the 'thinking' state, the user interacts with the AI model, describing the task they want to execute. The AI
//     responds with a concise explanation and a proposed bash command that corresponds to the described task.
//
//   - In the 'executing' state, the user reviews the proposed bash code, potentially editing it for correctness or
//     to better suit their specific needs. If the user approves the code (by typing "ok"), it is immediately executed in the shell.
//     If the user is not satisfied with the proposed code, they can continue interacting with the AI to improve it.
//
// The transition between states is managed by the REPL (Read-Eval-Print Loop) object, which represents the main
// interface of the utility. It maintains the state of the utility and routes user input to the appropriate
// component (the AI model or the shell executor).
//
// The utility also supports other features like command history, which is stored in a .think_history file in the
// user's home directory, and configurable prompts for different states, which are color-coded for better visibility.
//
// This utility uses several third-party packages, including 'github.com/chzyer/readline' for readline capabilities,
// 'github.com/fatih/color' for colorizing output, 'github.com/sirupsen/logrus' for logging, and 'github.com/ryszard/agency'
// for interacting with the AI model.
//
// Usage:
//
// The utility is a command-line tool, and as such is run from the terminal. You can set the AI model to use with the
// '-model' flag, or by setting the 'THINK_MODEL' environment variable. If neither is provided, the default model is 'gpt-4'.
//
// You can also provide an initial command to the utility by passing it as an argument.
//
// Example:
//
//	think -model=gpt-4 "list all files in this directory"
package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/ryszard/agency/agent"
	"github.com/ryszard/agency/client"
	"github.com/ryszard/agency/client/openai"
	"github.com/sirupsen/logrus"
)

var (
	model = flag.String("model", "", "model to use. The default is gpt-4")
)

//go:embed system.md
var SystemPrompt string

//go:embed user.md
var UserPrompt string

type REPL struct {
	readline        *readline.Instance
	inCode          bool
	agent           agent.Agent
	shellPath       string
	thinkingPrompt  string
	executingPrompt string
}

func NewREPL(agent agent.Agent, shellPath, initialInput string) *REPL {
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgCyan).SprintFunc()
	repl := &REPL{
		agent:           agent,
		shellPath:       shellPath,
		thinkingPrompt:  fmt.Sprintf("%s> ", blue("think")),
		executingPrompt: fmt.Sprintf("%s> ", red("run")),
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	historyFile := filepath.Join(homeDir, ".think_history")

	config := &readline.Config{
		Prompt:            "> ",
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		panic(err)
	}

	rl.WriteStdin([]byte(initialInput))

	repl.readline = rl
	return repl
}

func (repl *REPL) Close() {
	repl.readline.Close()
}

func (repl *REPL) intoCodeLoop() {
	repl.inCode = true
	repl.readline.SetPrompt(repl.executingPrompt)
}

func (repl *REPL) outOfCodeLoop() {
	repl.inCode = false
	repl.readline.SetPrompt(repl.thinkingPrompt)
}

func (repl *REPL) Run() {
	repl.outOfCodeLoop()

	var lastOut, lastErr string
	var exitCode int
	var commandRun bool
	for {
		line, err := repl.readline.Readline()
		if err != nil { // io.EOF
			if err == readline.ErrInterrupt {
				if repl.inCode {
					repl.outOfCodeLoop()
				}
				continue
			} else if err == io.EOF {
				if repl.inCode {
					repl.outOfCodeLoop()
					continue
				}
				return
			} else {
				log.Fatal(err)
			}
		}
		if repl.inCode {
			if line == "" {
				repl.outOfCodeLoop()
				continue
			}
			var stdoutBuf, stderrBuf bytes.Buffer

			cmd := exec.Command(repl.shellPath, "-c", line)

			// Create multiwriters so we write to both the buffers and standard output/error
			outMulti := io.MultiWriter(os.Stdout, &stdoutBuf)
			errMulti := io.MultiWriter(os.Stderr, &stderrBuf)

			cmd.Stdout = outMulti
			cmd.Stderr = errMulti
			cmd.Run()
			lastOut, lastErr = stdoutBuf.String(), stderrBuf.String()
			// truncate lastOut and lastErr to 1000 characters
			if len(lastOut) > 1000 {
				lastOut = lastOut[len(lastOut)-1000:]
			}
			if len(lastErr) > 1000 {
				lastErr = lastErr[len(lastErr)-1000:]
			}
			commandRun = true
			exitCode = cmd.ProcessState.ExitCode()

			repl.outOfCodeLoop()
		} else {
			if _, err := repl.agent.Listen("user", struct {
				Message    string
				CommandRun bool
				Stdout     string
				Stderr     string
				ExitCode   int
			}{
				Message:    strings.TrimSpace(line),
				CommandRun: commandRun,
				Stdout:     lastOut,
				Stderr:     lastErr,
				ExitCode:   exitCode,
			}); err != nil {
				log.Fatal(err)
			}

			response, err := repl.agent.Respond(context.Background(), agent.WithStreaming(os.Stdout))
			if err != nil {
				log.Fatal(err)
			}
			respLines := strings.Split(response, "\n")
			// remove empty strings from respLines
			for i := len(respLines) - 1; i >= 0; i-- {
				if respLines[i] == "" {
					respLines = append(respLines[:i], respLines[i+1:]...)
				}
			}
			command := respLines[len(respLines)-1]
			if _, err := repl.readline.WriteStdin([]byte(command)); err != nil {
				log.Fatal(err)
			}

			repl.intoCodeLoop()
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `think is a command-line tool that uses AI to generate and execute bash commands.

Usage:

  think [-model model] "your command"

Options:

  -model    Specifies the AI model to use. The default model is 'gpt-4'. 

You can also set the AI model to use with the 'THINK_MODEL' environment variable. If neither is provided, 
the default model is 'gpt-4'. 

Examples:

  think "list all files in this directory"
  think -model=gpt-4 "create a new directory called test"

`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	flag.Parse()

	logrus.SetLevel(logrus.ErrorLevel)
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

	ag := agent.New("scripter", agent.WithClient(cl), agent.WithModel(*model), agent.WithMaxTokens(500))
	ag, err := agent.Templated(ag, map[string]string{
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

	repl := NewREPL(ag, shellPath, strings.Join(flag.Args(), " "))
	defer repl.Close()

	repl.Run()

}
