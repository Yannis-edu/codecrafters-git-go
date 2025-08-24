package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// catFileCmd represents the cat-file command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide contents or details of repository objects",
	Args:  cobra.ExactArgs(1),
	Run:   catFile,
}

func init() {
	rootCmd.AddCommand(catFileCmd)
	catFileCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the contents of <object> based on its type.") // Not used yet
}

func catFile(cmd *cobra.Command, args []string) {
	filename := args[0]
	filename = ".git/objects/" + filename[:2] + "/" + filename[2:]

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	raw, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	parts := bytes.SplitN(raw, []byte{0}, 2)
	fmt.Print(string(parts[1]))
}
