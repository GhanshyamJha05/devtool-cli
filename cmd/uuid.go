package cmd

import (
	"crypto/rand"
	"fmt"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a random UUID v4",
	Long: `Generate a random UUID (version 4).

Example:
  devtool uuid`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uuid, err := generateUUID()
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to generate UUID: %v", err))
			return err
		}

		fmt.Println(uuid)
		return nil
	},
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Set the version to 4
	b[6] = (b[6] & 0x0f) | 0x40
	// Set the variant to RFC4122
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func init() {
	rootCmd.AddCommand(uuidCmd)
}
