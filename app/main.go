package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func getBuiltins() map[string]struct{} {
	return map[string]struct{}{
		"exit": {},
		"echo": {},
		"type": {},
	}
}

func handleExit(commandArgs []string) {
	if len(commandArgs) != 1 {
		fmt.Fprintf(os.Stderr, "exit command takes one argument")
		return
	}
	switch commandArgs[0] {
	case "0":
		os.Exit(0)
	case "1":
		os.Exit(1)
	}
}

func handleEcho(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Fprintf(os.Stderr, "echo command takes at least one argument")
		return
	}

	var builder strings.Builder
	for _, arg := range commandArgs {
		builder.WriteString(arg + " ")
	}
	builder.WriteString("\n")

	fmt.Fprint(os.Stdout, builder.String())
}

func handleType(commandArgs []string) {
	if len(commandArgs) != 1 {
		fmt.Fprintf(os.Stderr, "type command takes one argument")
		return
	}

	commandName := commandArgs[0]
	builtins := getBuiltins()
	if _, ok := builtins[commandName]; ok {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", commandName)
		return
	}

	if binPath, found := findBin(commandName); found {
		fmt.Fprintf(os.Stdout, "%s is %s\n", commandName, binPath)
		return
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", commandName)
}

func isExecutable(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}

	return fi.Mode()&0o111 != 0
}

func findBin(binName string) (string, bool) {
	paths := os.Getenv("PATH")
	for path := range strings.SplitSeq(paths, ":") {
		binPath := path + "/" + binName
		if fileInfo, err := os.Stat(binPath); err == nil {
			if !isExecutable(fileInfo) {
				continue
			}

			return binPath, true
		}
	}

	return "", false
}

func commandHandler(command string, commandArgs []string) {
	switch command {
	case "echo":
		handleEcho(commandArgs)
	case "exit":
		handleExit(commandArgs)
	case "type":
		handleType(commandArgs)

	default:
		fmt.Fprintln(os.Stdout, command+": not found")
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sigs
		fmt.Fprintf(os.Stdout, "\nreceived signal %s, exiting\n", s)
		os.Exit(0)
	}()

	for {
		fmt.Fprint(os.Stdout, "$ ")

		commandStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read input: %v\n", err)
			os.Exit(1)
		}

		commandStr = strings.TrimRight(commandStr, "\r\n")
		tokens := strings.Split(commandStr, " ")
		if len(tokens) == 0 {
			continue
		}

		command := tokens[0]
		commandArgs := tokens[1:]
		commandHandler(command, commandArgs)
	}
}
