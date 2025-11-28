package main

import (
	"fmt"
	"os"
	"strings"
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

	for _, r := range input {
		switch quoteChar {
		case 0:
			switch r {
			case '\'':
				if !escaped {
					quoteChar = '\''
				} else {
					sb.WriteRune(r)
				}
			case '"':
				if !escaped {
					quoteChar = '"'
				} else {
					sb.WriteRune(r)
				}
			case '\\':
				escaped = true
			case ' ':
				if escaped {
					sb.WriteRune(r)
					escaped = false
					break
				}

				str := strings.TrimSpace(sb.String())
				if len(str) != 0 {
					args = append(args, strings.Trim(sb.String(), " "))
				}

				sb.Reset()
			default:
				sb.WriteRune(r)
			}
		case '\'':
			switch r {
			case '\'':
				quoteChar = 0
			default:
				sb.WriteRune(r)
			}
		case '"':
			switch r {
			case '\\':
				if !escaped {
					escaped = true
				} else {
					sb.WriteRune(r)
					escaped = false
				}
			case '"':
				if !escaped {
					quoteChar = 0
				} else {
					sb.WriteRune(r)
					escaped = false
				}
			default:
				if escaped {
					sb.WriteRune('\\')
					escaped = false
				}
				sb.WriteRune(r)
			}
		}
	}

	if quoteChar != 0 {
		return nil, fmt.Errorf("unclosed quote in argument string")
	}

	if sb.Len() > 0 {
		args = append(args, sb.String())
	}

	return args, nil
}
