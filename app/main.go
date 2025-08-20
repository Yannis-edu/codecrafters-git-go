package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		filename := os.Args[3] // TODO: Use real a CLI lib tool
		filename = ".git/objects/" + filename[:2] + "/" + filename[2:]

		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			panic(err)
		}

		r, err := zlib.NewReader(f)
		defer r.Close()
		if err != nil {
			panic(err)
		}

		raw, err := io.ReadAll(r)
		if err != nil {
			panic(err)
		}

		parts := bytes.SplitN(raw, []byte{0}, 2)
		fmt.Print(string(parts[1]))

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
