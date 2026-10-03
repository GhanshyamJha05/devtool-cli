package cmd

import (
	"fmt"
	"net/http"

	"github.com/GhanshyamJha05/devtool-cli/internal/utils"

	"github.com/spf13/cobra"
)

var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local HTTP server",
	Long: `Start a simple local HTTP server to serve static files from the current directory.
	
Example:
  devtool serve
  devtool serve --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := ":" + port
		utils.PrintSuccess(fmt.Sprintf("Serving HTTP on 0.0.0.0 port %s (http://localhost:%s/)", port, port))
		
		http.Handle("/", http.FileServer(http.Dir(".")))
		
		if err := http.ListenAndServe(addr, nil); err != nil {
			utils.PrintError(fmt.Sprintf("Server failed: %v", err))
			return err
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to listen on")
}
