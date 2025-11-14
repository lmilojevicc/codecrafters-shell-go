package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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

func handleExit(args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "exit command takes one argument")
		return
	}
	switch args[0] {
	case "0":
		os.Exit(0)
	case "1":
		os.Exit(1)
	}
}

func handleEcho(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "echo command takes at least one argument")
		return
	}

	fmt.Fprintln(os.Stdout, strings.Join(args, " "))
}

func handleType(args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "type command takes one argument")
		return
	}

	commandName := args[0]
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

func handleCd(args []string) {
	if len(args) != 1 && len(args) != 0 {
		fmt.Fprintln(os.Stderr, "cd command takes one or zero arguments")
		return
	}

	path := args[0]
	if !pathExists(path) {
		return
	}

	os.Chdir(path)
}

func pathExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
		return false
	}

	if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not directory\n", path)
		return false
	}

	return true
}

func executeBinary(bin string, args []string) {
	args = append([]string{bin}, args...)

	if binPath, found := findBin(args[0]); found {
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
	case "pwd":
		hadnlePwd()
	case "cd":
		handleCd(commandArgs)
	default:
		executeBinary(command, commandArgs)
	}
}

func hadnlePwd() {
	currPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't get current working directory:", err)
	}

	fmt.Fprintln(os.Stdin, currPath)
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

		tokens := strings.Fields(input)
		if len(tokens) == 0 {
			continue
		}

		command := tokens[0]
		commandArgs := tokens[1:]
		commandHandler(command, commandArgs)
	}
}
