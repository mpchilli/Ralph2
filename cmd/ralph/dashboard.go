package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ralph2/internal/dashboard"
	"ralph2/pkg/utils"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start the Web Dashboard Server",
	Run: func(cmd *cobra.Command, args []string) {
		bus := utils.NewEventBus()
		server := dashboard.NewServer(bus)
		
		fmt.Println("🌐 Starting Dashboard at http://localhost:8080")
		if err := server.Start(":8080"); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
