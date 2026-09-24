package cmd

import (
	"fmt"
	"os"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var saveFormatted string
var forceSaveFormatted bool

var formatCmd = &cobra.Command{
	Use:   "format <file.json>",
	Short: "Pretty-print JSON files",
	Long: `Read a JSON file, validate its structure, and output a beautifully formatted version.

Examples:
  devtool format config.json
  devtool format data.json --save pretty.json
  devtool format response.json --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		if Verbose {
			utils.PrintDebug(fmt.Sprintf("Input file: %s", filePath))
		}

		utils.PrintInfo(fmt.Sprintf("Formatting %s...", filePath))

		formatted, err := utils.FormatJSON(filePath)
		if err != nil {
			utils.PrintError(err.Error())
			return err
		}

		// Save to file or print to stdout
		if saveFormatted != "" {
			if !forceSaveFormatted {
				if _, err := os.Stat(saveFormatted); err == nil {
					utils.PrintError(fmt.Sprintf("File %s already exists. Use --force to overwrite.", saveFormatted))
					return fmt.Errorf("file exists")
				}
			}
			if Verbose {
				utils.PrintDebug(fmt.Sprintf("Writing output to: %s", saveFormatted))
			}
			if err := os.WriteFile(saveFormatted, []byte(formatted+"\n"), 0644); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to save: %s", err.Error()))
				return err
			}
			utils.PrintSuccess(fmt.Sprintf("Formatted JSON saved to %s", saveFormatted))
			return nil
		}

		fmt.Println()
		fmt.Println(formatted)
		fmt.Println()
		utils.PrintSuccess("JSON is valid and formatted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
	formatCmd.Flags().StringVarP(&saveFormatted, "save", "s", "", "save formatted output to a file")
	formatCmd.Flags().BoolVarP(&forceSaveFormatted, "force", "f", false, "force overwrite if file already exists")
}
