package cmd

import (
	"context"
	"fmt"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update devtool-cli to the latest version",
	Long:  `Check for a newer version of devtool-cli on GitHub and update it automatically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner, _ := pterm.DefaultSpinner.Start("Checking for updates...")
		
		// The selfupdate library requires a context
		ctx := context.Background()
		
		updater, err := selfupdate.NewUpdater(selfupdate.Config{
			Validator: &selfupdate.ChecksumValidator{
				UniqueFilename: "checksums.txt",
			},
		})
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to initialize updater: %v", err))
			return err
		}

		latest, err := updater.UpdateSelf(ctx, version, selfupdate.ParseSlug("GhanshyamJha05/devtool-cli"))
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to update: %v", err))
			return err
		}

		if latest.Version() == version {
			spinner.Success(fmt.Sprintf("You are already on the latest version (v%s)", version))
			return nil
		}

		spinner.Success(fmt.Sprintf("Successfully updated to %s", latest.Version()))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
