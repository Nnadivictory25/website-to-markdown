package main

import (
	"fmt"
	"os"

	"website-markdown/cmd"
	"website-markdown/internal/api"

	"github.com/spf13/cobra"
)

var (
	serverMode bool
	port       string
)

var rootCmd = &cobra.Command{
	Use:   "website-markdown",
	Short: "🔄 Convert websites to markdown - CLI and API server",
	Long: `Website to Markdown Converter

This tool can run in two modes:
1. CLI mode (default): Convert websites directly from command line
2. Server mode: Run as HTTP API server for web applications

Examples:
  # CLI mode
  website-markdown https://example.com
  website-markdown https://example.com --depth 2 --output ./docs

  # Server mode  
  website-markdown --server --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverMode {
			fmt.Println("🚀 Starting in server mode...")
			return api.StartAPIServer(port)
		}

		// If no URL provided and not in server mode, show help
		if len(args) == 0 {
			return cmd.Help()
		}

		// Delegate to CLI command
		return cmd.Execute()
	},
}

func init() {
	rootCmd.Flags().BoolVar(&serverMode, "server", false, "Run as API server")
	rootCmd.Flags().StringVar(&port, "port", "8080", "API server port (only used with --server)")
}

func main() {
	if len(os.Args) > 1 {
		// Check if running in server mode
		if os.Args[1] == "--server" || os.Args[1] == "-s" {
			fmt.Println("🚀 Starting API server...")
			if err := api.StartAPIServer(port); err != nil {
				fmt.Printf("❌ Server failed to start: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Check for help flags
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			rootCmd.Help()
			return
		}
	}

	// Default to CLI mode
	if err := cmd.Execute(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
