package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sigs
		fmt.Fprintf(os.Stdout, "received signal %s, exiting\n", s)
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
		commandWithArgs := strings.Split(commandStr, " ")

		command := commandWithArgs[0]
		commandArgs := commandWithArgs[1:]

		switch command {
		case "exit":
			if len(commandArgs) < 1 {
				fmt.Println("exit command takes one argument")
				continue
			}

			switch commandArgs[0] {
			case "0":
				os.Exit(0)
			case "1":
				os.Exit(1)
			}
		}

		fmt.Fprintln(os.Stdout, command+": not found")
	}
}
