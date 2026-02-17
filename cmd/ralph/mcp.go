package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ralph2/internal/mcp"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the Ralph2 MCP Server (stdio)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "🤖 Starting Ralph2 MCP Server...\n")
		
		s := mcp.NewMCPServer()
		if err := s.Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ MCP Server failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
