package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getBuiltins() map[string]struct{} {
	return map[string]struct{}{
		"cd":   {},
		"echo": {},
		"exit": {},
		"pwd":  {},
		"type": {},
	}
}

func HandleExit(args []string) {
	if len(args) == 0 {
		os.Exit(0)
	}

	switch args[0] {
	case "0":
		os.Exit(0)
	case "1":
		os.Exit(1)
	default:
		fmt.Fprintln(os.Stderr, "invalid argument for exit command")
		return
	}
}

func HandleEcho(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "echo command takes at least one argument")
		return
	}

	fmt.Fprintln(os.Stdout, strings.Join(args, " "))
}

func HandleType(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "type command takes one argument")
		return
	}

	commandName := args[0]
	builtins := getBuiltins()
	if _, ok := builtins[commandName]; ok {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", commandName)
		return
	}

	if binPath, found := FindBin(commandName); found {
		fmt.Fprintf(os.Stdout, "%s is %s\n", commandName, binPath)
		return
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", commandName)
}

func HandleCd(args []string) {
	if len(args) != 1 && len(args) != 0 {
		fmt.Fprintln(os.Stderr, "cd command takes one or zero arguments")
		return
	}

	path := args[0]
	isHome := path[0] == '~'
	if isHome {
		home := os.Getenv("HOME")
		path = strings.Replace(path, "~", home, 1)
	}

	if !PathExists(path) {
		return
	}

	os.Chdir(path)
}

func HandlePwd() {
	currPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't get current working directory:", err)
	}

	fmt.Fprintln(os.Stdin, currPath)
}

func ExecuteBinary(bin string, args []string) {
	args = append([]string{bin}, args...)

	if binPath, found := FindBin(args[0]); found {
		cmd := exec.Command(binPath, args[1:]...)
		cmd.Args = args
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stdout, "Error executing command:", err)
		}

		return
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", bin)
}
