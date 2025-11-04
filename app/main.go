package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read input: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, command[:len(command)-1]+": not found")
}
