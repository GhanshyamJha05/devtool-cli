package cmd

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var useMD5 bool
var useSHA256 bool
var isString bool

var hashCmd = &cobra.Command{
	Use:   "hash <file_or_string>",
	Short: "Compute MD5 or SHA-256 hash",
	Long: `Compute the MD5 or SHA-256 hash of a file (default) or a direct string.
By default, it computes SHA-256.

Example:
  devtool hash config.json
  devtool hash --md5 config.json
  devtool hash --string "my secret string"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		
		var data []byte
		
		if isString {
			data = []byte(target)
		} else {
			file, err := os.Open(target)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to open file: %v", err))
				return err
			}
			defer file.Close()
			
			d, err := io.ReadAll(file)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to read file: %v", err))
				return err
			}
			data = d
		}
		
		if useMD5 {
			hash := md5.Sum(data)
			fmt.Printf("%x\n", hash)
		} else {
			// default to SHA256 if nothing or only sha256 is specified
			hash := sha256.Sum256(data)
			fmt.Printf("%x\n", hash)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().BoolVar(&useMD5, "md5", false, "compute MD5 hash instead of SHA-256")
	hashCmd.Flags().BoolVar(&useSHA256, "sha256", true, "compute SHA-256 hash (default)")
	hashCmd.Flags().BoolVar(&isString, "string", false, "treat the argument as a string rather than a filepath")
}
