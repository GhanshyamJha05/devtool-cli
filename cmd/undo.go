package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo <folder>",
	Short: "Undo a previous clean operation",
	Long: `Revert the files back to their original locations using the .devtool-undo.json mapping file left in the directory after a clean operation.

Example:
  devtool undo ./downloads`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		folderPath := args[0]
		undoPath := filepath.Join(folderPath, ".devtool-undo.json")

		if _, err := os.Stat(undoPath); os.IsNotExist(err) {
			utils.PrintError(fmt.Sprintf("Undo file not found: %s", undoPath))
			return fmt.Errorf("undo file not found")
		}

		moves, err := utils.LoadUndoFile(undoPath)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to load undo file: %v", err))
			return err
		}

		utils.PrintInfo(fmt.Sprintf("Undoing %d file moves...", len(moves)))

		var errCount int
		for _, move := range moves {
			// move.NewPath is where it is now, move.OldPath is where it was
			if Verbose {
				utils.PrintDebug(fmt.Sprintf("Moving %s -> %s", move.NewPath, move.OldPath))
			}
			if err := os.Rename(move.NewPath, move.OldPath); err != nil {
				utils.PrintWarning(fmt.Sprintf("Failed to move %s: %v", move.NewPath, err))
				errCount++
			}
		}

		// Delete the undo file if everything succeeded
		if errCount == 0 {
			os.Remove(undoPath)
			utils.PrintSuccess("Undo complete! All files restored successfully.")
		} else {
			utils.PrintWarning(fmt.Sprintf("Undo finished with %d errors. The undo file was kept.", errCount))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
