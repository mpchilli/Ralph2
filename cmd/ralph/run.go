package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ralph2/internal/service"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start an autonomous coding task",
	Long:  `Run a full autonomous cycle (Planning, Building, Verifying) based on a prompt.`,
	Run: func(cmd *cobra.Command, args []string) {
		if prompt == "" && len(args) > 0 {
			prompt = args[0]
		}
		if prompt == "" {
			fmt.Println("Error: prompt is required")
			os.Exit(1)
		}
		
		svc := service.NewOrchestratorService()
		if err := svc.Run(prompt, complexity); err != nil {
			fmt.Printf("❌ Task failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "User request prompt")
	runCmd.Flags().StringVarP(&complexity, "complexity", "c", "streamlined", "Loop complexity (fast, streamlined, full)")
	rootCmd.AddCommand(runCmd)
}
