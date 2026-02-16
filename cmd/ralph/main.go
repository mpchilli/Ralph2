package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ralph2/internal/core"
	"ralph2/pkg/utils"
)

var (
	complexity string
	prompt     string
)

var rootCmd = &cobra.Command{
	Use:   "ralph2",
	Short: "Ralph 2.1 - Autonomous Dev Orchestrator",
	Long:  `Ralph 2.1 runs in Loki Mode - Autonomous Planning, Coding, and Verifying.`,
	Run: func(cmd *cobra.Command, args []string) {
		runOrchestrator()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&complexity, "complexity", "c", "streamlined", "Complexity mode: fast, streamlined, full")
	rootCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Initial prompt for the task")
	rootCmd.AddCommand(tuiCmd)
}

func runOrchestrator() {
	fmt.Printf("Initializing Ralph 2.1 with Complexity: %s\n", complexity)
	
	if prompt != "" {
		// Init EventBus
		bus := utils.NewEventBus()
		
		// Init FSM with EventBus
		fsm := core.NewStateManager(bus)
		fsm.TransitionTo(core.StatePlanning)
		
		fmt.Printf("Received Prompt: %q\n", prompt)
		// Trigger One-Shot Logic here
	} else {
		fmt.Println("No prompt provided. Entering Interactive Mode (TUI)...")
		startTUI()
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
