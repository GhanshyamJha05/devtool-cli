package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var saveFormatted string
var forceSaveFormatted bool
var inPlace bool

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
		targetPath := args[0]

		if Verbose {
			utils.PrintDebug(fmt.Sprintf("Input path: %s", targetPath))
		}

		info, err := os.Stat(targetPath)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Cannot access path: %v", err))
			return err
		}

		if info.IsDir() {
			if saveFormatted != "" {
				utils.PrintError("--save is not supported when formatting a directory")
				return fmt.Errorf("invalid arguments")
			}
			utils.PrintInfo(fmt.Sprintf("Formatting directory %s...", targetPath))
			var formatCount int
			err := filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && strings.ToLower(filepath.Ext(path)) == ".json" {
					if inPlace {
						tempFile := path + ".tmp"
						out, err := os.Create(tempFile)
						if err != nil {
							utils.PrintWarning(fmt.Sprintf("Failed to format %s: %v", path, err))
							return nil
						}
						err = utils.FormatJSON(path, out)
						out.Close()
						if err != nil {
							utils.PrintWarning(fmt.Sprintf("Failed to format %s: %v", path, err))
							os.Remove(tempFile)
						} else {
							os.Rename(tempFile, path)
							formatCount++
							if Verbose {
								utils.PrintDebug(fmt.Sprintf("Formatted %s", path))
							}
						}
					} else {
						fmt.Printf("--- %s ---\n", path)
						utils.FormatJSON(path, os.Stdout)
						fmt.Println()
						formatCount++
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
			utils.PrintSuccess(fmt.Sprintf("Processed %d JSON file(s)", formatCount))
			return nil
		}

		// Single file processing
		filePath := targetPath
		utils.PrintInfo(fmt.Sprintf("Formatting %s...", filePath))

		var out io.Writer
		var outFile *os.File

		if inPlace {
			tempFile := filePath + ".tmp"
			outFile, err = os.Create(tempFile)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to create temp file: %s", err.Error()))
				return err
			}
			defer outFile.Close()
			out = outFile
		} else if saveFormatted != "" {
			if !forceSaveFormatted {
				if _, err := os.Stat(saveFormatted); err == nil {
					utils.PrintError(fmt.Sprintf("File %s already exists. Use --force to overwrite.", saveFormatted))
					return fmt.Errorf("file exists")
				}
			}
			if Verbose {
				utils.PrintDebug(fmt.Sprintf("Writing output to: %s", saveFormatted))
			}
			
			var err error
			outFile, err = os.Create(saveFormatted)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to create file: %s", err.Error()))
				return err
			}
			defer outFile.Close()
			out = outFile
		} else {
			out = os.Stdout
			fmt.Println()
		}

		err = utils.FormatJSON(filePath, out)
		if err != nil {
			utils.PrintError(err.Error())
			return err
		}

		if inPlace {
			outFile.Close()
			if err := os.Rename(filePath+".tmp", filePath); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to overwrite file in place: %v", err))
				return err
			}
			utils.PrintSuccess(fmt.Sprintf("Formatted JSON in-place: %s", filePath))
		} else if saveFormatted != "" {
			utils.PrintSuccess(fmt.Sprintf("Formatted JSON saved to %s", saveFormatted))
		} else {
			fmt.Println()
			utils.PrintSuccess("JSON is valid and formatted")
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
	formatCmd.Flags().StringVarP(&saveFormatted, "save", "s", "", "save formatted output to a file")
	formatCmd.Flags().BoolVarP(&forceSaveFormatted, "force", "f", false, "force overwrite if file already exists")
	formatCmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "overwrite the original file(s) with formatted output")
}
