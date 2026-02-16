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

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the Terminal User Interface",
	Long:  `Run Ralph 2.1 in interactive TUI mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTUI()
	},
}

func startTUI() {
	// Init EventBus
	bus := utils.NewEventBus()

	// Init FSM with EventBus
	fsm := core.NewStateManager(bus)
	// Transition to initial state
	fsm.TransitionTo(core.StatePlanning)

	// Init TUI Model with EventBus
	model := tui.NewModel(bus, fsm.GetState())
	
	// Create and run program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
