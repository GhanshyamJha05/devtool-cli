package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

var dryRun bool

var cleanCmd = &cobra.Command{
	Use:   "clean <folder>",
	Short: "Organize files in a directory by file type",
	Long: `Scan a directory and automatically sort files into categorized subfolders
based on their file extension (Images, Documents, Code, Videos, etc.)

Examples:
  devtool clean ./downloads
  devtool clean C:\Users\you\Desktop --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folderPath := args[0]

		if Verbose {
			utils.PrintDebug(fmt.Sprintf("Target directory: %s", folderPath))
		}

		var spinner *pterm.SpinnerPrinter
		if dryRun {
			spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("[DRY-RUN] Scanning folder: %s...", folderPath))
		} else {
			spinner, _ = pterm.DefaultSpinner.Start(fmt.Sprintf("Scanning folder: %s...", folderPath))
		}

		result, err := utils.OrganizeFolder(folderPath, dryRun)
		if err != nil {
			spinner.Fail(err.Error())
			return err
		}
		spinner.Success("Scan complete")

		// Handle edge case: nothing to organize
		if result.TotalFiles == 0 {
			utils.PrintWarning("No files found to organize")
			return nil
		}

		// Print a summary of what was moved
		fmt.Println()
		for category, files := range result.Moved {
			utils.PrintSuccess(fmt.Sprintf("%-12s → %d file(s)", category, len(files)))
			if Verbose {
				for _, f := range files {
					utils.PrintDebug(fmt.Sprintf("  moved: %s", f))
				}
			}
		}

		if result.Skipped > 0 {
			utils.PrintWarning(fmt.Sprintf("Skipped %d file(s) with no extension", result.Skipped))
		}

		if !dryRun && len(result.Moves) > 0 {
			undoPath := filepath.Join(folderPath, ".devtool-undo.json")
			if err := utils.SaveUndoFile(undoPath, result.Moves); err != nil {
				utils.PrintWarning(fmt.Sprintf("Failed to save undo file: %v", err))
			} else {
				utils.PrintSuccess(fmt.Sprintf("Undo mapping saved to %s", undoPath))
			}
		}

		fmt.Println()
		if dryRun {
			utils.PrintSuccess(fmt.Sprintf("[DRY-RUN] Would organize %d file(s) into %d categories", result.TotalFiles, len(result.Moved)))
		} else {
			utils.PrintSuccess(fmt.Sprintf("Done! Organized %d file(s) into %d categories", result.TotalFiles, len(result.Moved)))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "simulate the organization without moving files")
}
