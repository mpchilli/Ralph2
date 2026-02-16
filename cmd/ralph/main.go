package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"ralph2/internal/core"
	"ralph2/internal/tui"
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
}

func runOrchestrator() {
	fmt.Printf("Initializing Ralph 2.1 with Complexity: %s\n", complexity)
	
	// Init EventBus
	bus := utils.NewEventBus()
	
	// Init FSM
	fsm := core.NewStateManager()
	fsm.TransitionTo(core.StatePlanning)
	
	if prompt != "" {
		fmt.Printf("Received Prompt: %q\n", prompt)
		// Trigger One-Shot Logic here
	} else {
		fmt.Println("No prompt provided. Entering Interactive Mode (TUI)...")
		
		// Init TUI
		model := tui.NewModel(bus)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
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
