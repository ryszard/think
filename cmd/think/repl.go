package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/ryszard/agency/agent"
	"github.com/sirupsen/logrus"
)

type REPL struct {
	readline        *readline.Instance
	inCode          bool
	agent           agent.Agent
	shellPath       string
	thinkingPrompt  string
	executingPrompt string
	sendOutput      bool
}

func NewREPL(agent agent.Agent, shellPath, initialInput string, sendOutput bool) *REPL {
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgCyan).SprintFunc()
	repl := &REPL{
		agent:           agent,
		shellPath:       shellPath,
		thinkingPrompt:  fmt.Sprintf("%s> ", blue("think")),
		executingPrompt: fmt.Sprintf("%s> ", red("run")),
		sendOutput:      sendOutput,
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
		AutoComplete:      new(FileCompleter),
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
	var commandWasRun bool
	var actualCommand string
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
			actualCommand = strings.Join([]string{repl.shellPath, "-c", line}, " ")
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
			commandWasRun = true
			exitCode = cmd.ProcessState.ExitCode()

			repl.outOfCodeLoop()
		} else {
			feedback, err := repl.agent.Listen("user", struct {
				Message       string
				CommandWasRun bool
				ActualCommand string
				Stdout        string
				Stderr        string
				ExitCode      int
				SendOutput    bool
			}{
				Message:       strings.TrimSpace(line),
				CommandWasRun: commandWasRun,
				ActualCommand: actualCommand,
				Stdout:        lastOut,
				Stderr:        lastErr,
				ExitCode:      exitCode,
				SendOutput:    repl.sendOutput,
			})
			if err != nil {
				log.Fatal(err)
			}
			logrus.WithField("feedback", feedback).Debug("feedback sent to the AI model")

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
