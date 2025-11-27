package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func PathExists(path string) bool {
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

func isExecutable(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}

	return fi.Mode()&0o111 != 0
}

func FindBin(binName string) (string, bool) {
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

func ProcessArgs(input string) ([]string, error) {
	if len(input) == 0 {
		return []string{}, nil
	}

	var args []string
	var sb strings.Builder

	var quoteChar rune
	escaped := false
	argStrated := false

	for _, r := range input {
		if escaped {
			escaped = false
			sb.WriteRune(r)
			continue
		}

		if quoteChar != 0 {
			if r == '\\' {
				escaped = true
				continue
			}
			if r == quoteChar {
				quoteChar = 0
			} else {
				sb.WriteRune(r)
			}

			continue
		}

		switch {
		case unicode.IsSpace(r):
			if sb.Len() == 0 {
				continue
			}
			args = append(args, sb.String())
			sb.Reset()
			argStrated = false
		case r == '\'' || r == '"':
			quoteChar = r
			argStrated = true
		case r == '\\':
			sb.WriteRune(r)
			argStrated = true
		default:
			sb.WriteRune(r)
			argStrated = true
		}
	}

	if quoteChar != 0 {
		return nil, fmt.Errorf("unclosed quote in argument string")
	}

	if argStrated {
		args = append(args, sb.String())
	}

	return args, nil
}
