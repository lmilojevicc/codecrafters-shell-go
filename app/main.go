package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func commandHandler(command string, args []string) {
	switch command {
	case "echo":
		HandleEcho(args)
	case "exit":
		HandleExit(args)
	case "type":
		HandleType(args)
	case "pwd":
		HandlePwd()
	case "cd":
		HandleCd(args)
	default:
		ExecuteBinary(command, args)
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

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read input: %v\n", err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		tokens, _ := ProcessArgs(input)

		command := tokens[0]
		commadnArgs := tokens[1:]

		commandHandler(command, commadnArgs)
	}
}
