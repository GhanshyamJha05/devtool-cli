package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time <timestamp|now>",
	Short: "Convert Unix timestamps to human-readable dates",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		
		if input == "now" {
			now := time.Now()
			pterm.Success.Printf("Current Unix Epoch: %d\n", now.Unix())
			pterm.Info.Printf("Local Time: %s\n", now.Format(time.RFC1123))
			pterm.Info.Printf("UTC Time:   %s\n", now.UTC().Format(time.RFC1123))
			return nil
		}
		
		// Attempt to parse as int64
		epoch, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			utils.PrintError("Invalid timestamp format. Provide a valid Unix epoch (e.g. 1700000000) or 'now'.")
			return err
		}
		
		// Determine if it's seconds, milliseconds, etc.
		// A standard epoch in seconds is around 10 digits.
		// If it's 13 digits, it's milliseconds.
		var t time.Time
		if epoch > 9999999999 {
			// Assume milliseconds
			t = time.UnixMilli(epoch)
			pterm.Info.Println("Detected format: Milliseconds")
		} else {
			t = time.Unix(epoch, 0)
			pterm.Info.Println("Detected format: Seconds")
		}
		
		fmt.Println()
		pterm.Success.Printf("Local Time: %s\n", t.Format(time.RFC1123))
		pterm.Success.Printf("UTC Time:   %s\n", t.UTC().Format(time.RFC1123))
		
		// Tell relative time
		if t.Before(time.Now()) {
			pterm.DefaultBasicText.Printf("(%s ago)\n", time.Since(t).Round(time.Second))
		} else {
			pterm.DefaultBasicText.Printf("(in %s)\n", time.Until(t).Round(time.Second))
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
