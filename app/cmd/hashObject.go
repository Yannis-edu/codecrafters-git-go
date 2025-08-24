package cmd

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hash-object command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Compute object ID and optionally create an object from a file",
	Args:  cobra.ExactArgs(1),
	RunE:  hashObject,
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().BoolP("write", "w", false, "Actually write the object into the object database.")
}

func hashObject(cmd *cobra.Command, args []string) error {
	hasToWrite, err := cmd.Flags().GetBool("write")
	if err != nil {
		return err
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	data := append([]byte("blob "), []byte(strconv.Itoa(len(content)))...)
	data = append(data, byte(0))
	data = append(data, content...)

	h := sha1.New()
	h.Write(data)
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println(sha)

	if hasToWrite {
		dir := ".git/objects/" + sha[:2] + "/"

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		fout, err := os.Create(dir + sha[2:])
		if err != nil {
			return err
		}
		defer fout.Close()

		zw := zlib.NewWriter(fout)
		defer zw.Close()
		_, err = zw.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
