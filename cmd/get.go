package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var saveFile string
var forceSave bool

type progressReader struct {
	io.Reader
	bar *pterm.ProgressbarPrinter
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	if pr.bar != nil && n > 0 {
		pr.bar.Add(n)
	}
	return
}

var getCmd = &cobra.Command{
	Use:   "get <url>",
	Short: "Fetch API data and display response",
	Long: `Fetch data from any HTTP/HTTPS URL and display the response.
Supports saving output to a file and verbose debug mode.

Examples:
  devtool get https://jsonplaceholder.typicode.com/posts/1
  devtool get https://api.github.com/users/octocat --save user.json
  devtool get https://httpbin.org/get --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rawURL := args[0]

		if Verbose {
			utils.PrintDebug(fmt.Sprintf("Target URL: %s", rawURL))
		}

		utils.PrintInfo(fmt.Sprintf("Fetching data from %s...", rawURL))

		resp, err := utils.FetchData(rawURL)
		if err != nil {
			utils.PrintError(err.Error())
			return err
		}
		defer resp.BodyReader.Close()

		// Display response metadata
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			utils.PrintSuccess(fmt.Sprintf("Status: %s  |  Time: %s", resp.Status, resp.Duration.Round(1000000)))
		} else {
			utils.PrintWarning(fmt.Sprintf("Status: %s  |  Time: %s", resp.Status, resp.Duration.Round(1000000)))
		}

		if Verbose {
			utils.PrintDebug(fmt.Sprintf("Content-Type: %s", resp.Headers.Get("Content-Type")))
			utils.PrintDebug(fmt.Sprintf("Content-Length: %s", resp.Headers.Get("Content-Length")))
		}

		// Save to file or print to stdout
		if saveFile != "" {
			if !forceSave {
				if _, err := os.Stat(saveFile); err == nil {
					utils.PrintError(fmt.Sprintf("File %s already exists. Use --force to overwrite.", saveFile))
					return fmt.Errorf("file exists")
				}
			}

			file, err := os.Create(saveFile)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to create file: %s", err.Error()))
				return err
			}
			defer file.Close()

			var reader io.Reader = resp.BodyReader
			var bar *pterm.ProgressbarPrinter
			if resp.ContentLength > 0 {
				bar, _ = pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength)).WithTitle("Downloading").Start()
				reader = &progressReader{Reader: resp.BodyReader, bar: bar}
			}

			written, err := io.Copy(file, reader)
			if bar != nil {
				bar.Stop()
			}

			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to save: %s", err.Error()))
				return err
			}
			utils.PrintSuccess(fmt.Sprintf("Response saved to %s (%d bytes)", saveFile, written))
			return nil
		}

		fmt.Println()
		bodyBytes, err := io.ReadAll(resp.BodyReader)
		if err != nil {
			return err
		}
		fmt.Println(string(bodyBytes))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&saveFile, "save", "s", "", "save response to a file (e.g., --save output.json)")
	getCmd.Flags().BoolVarP(&forceSave, "force", "f", false, "force overwrite if file already exists")
}
