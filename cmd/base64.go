package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var decodeMode bool

var base64Cmd = &cobra.Command{
	Use:   "base64 <string>",
	Short: "Encode or decode Base64 strings",
	Long: `Encode a string to Base64 or decode a Base64 string back to plaintext.

Example:
  devtool base64 "hello world"
  devtool base64 --decode "aGVsbG8gd29ybGQ="`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		if decodeMode {
			decoded, err := base64.StdEncoding.DecodeString(target)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to decode base64: %v", err))
				return err
			}
			fmt.Println(string(decoded))
		} else {
			encoded := base64.StdEncoding.EncodeToString([]byte(target))
			fmt.Println(encoded)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().BoolVarP(&decodeMode, "decode", "d", false, "decode the input string instead of encoding")
}
