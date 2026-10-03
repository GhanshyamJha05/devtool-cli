package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt <token>",
	Short: "Inspect and decode JSON Web Tokens (JWT)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]
		parts := strings.Split(token, ".")

		if len(parts) != 3 {
			utils.PrintError("Invalid JWT format. Must contain 3 parts separated by dots.")
			return fmt.Errorf("invalid jwt")
		}

		// JWTs use Base64Url encoding without padding
		decode := func(s string) ([]byte, error) {
			if l := len(s) % 4; l > 0 {
				s += strings.Repeat("=", 4-l)
			}
			return base64.URLEncoding.DecodeString(s)
		}

		headerBytes, err := decode(parts[0])
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to decode header: %v", err))
			return err
		}

		payloadBytes, err := decode(parts[1])
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to decode payload: %v", err))
			return err
		}

		pterm.DefaultSection.Println("Header")
		if err := utils.FormatJSONBytes(headerBytes, os.Stdout); err != nil {
			fmt.Println(string(headerBytes)) // Fallback if not JSON
		}

		pterm.DefaultSection.Println("Payload")
		if err := utils.FormatJSONBytes(payloadBytes, os.Stdout); err != nil {
			fmt.Println(string(payloadBytes)) // Fallback
		}

		// Try to parse 'exp' claim to show expiration details
		var payload map[string]interface{}
		if err := json.Unmarshal(payloadBytes, &payload); err == nil {
			if expFloat, ok := payload["exp"].(float64); ok {
				expTime := time.Unix(int64(expFloat), 0)
				pterm.DefaultSection.Println("Expiration")
				if time.Now().After(expTime) {
					pterm.Error.Printf("EXPIRED: %v (%s ago)\n", expTime, time.Since(expTime).Round(time.Second))
				} else {
					pterm.Success.Printf("VALID: %v (expires in %s)\n", expTime, time.Until(expTime).Round(time.Second))
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)
}
