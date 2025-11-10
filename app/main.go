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
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read input: %v\n", err)
			os.Exit(1)
		}

		command = strings.TrimRight(command, "\r\n")
		fmt.Fprintln(os.Stdout, command+": not found")
	}
}
