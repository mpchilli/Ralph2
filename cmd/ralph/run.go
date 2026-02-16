package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ralph2/internal/core"
	"ralph2/internal/planner"
	"ralph2/pkg/utils"
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
		executeRun(prompt, complexity)
	},
}

func init() {
	runCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "User request prompt")
	runCmd.Flags().StringVarP(&complexity, "complexity", "c", "streamlined", "Loop complexity (fast, streamlined, full)")
	rootCmd.AddCommand(runCmd)
}

func executeRun(p string, c string) {
	fmt.Printf("🚀 Starting Ralph2 with prompt: %q\n", p)
	fmt.Printf("Complexity: %s\n", c)

	bus := utils.NewEventBus()
	sm := core.NewStateManager(bus)

	// Validate and set complexity
	level := core.ComplexityLevel(c)
	switch level {
	case core.ComplexityFast, core.ComplexityStreamlined, core.ComplexityFull:
		sm.SetComplexity(level)
	default:
		fmt.Printf("Warning: Unknown complexity %q, defaulting to streamlined\n", c)
	}

	// Move to planning
	sm.TransitionTo(core.StatePlanning)

	// Execute Planning
	pln := planner.NewHeuristicPlanner("spec.md")
	if err := pln.Plan(p); err != nil {
		fmt.Printf("Error during planning: %v\n", err)
		_ = sm.TransitionTo(core.StateFailed)
		return
	}

	fmt.Println("Plan generated successfully in spec.md")

	// Move to building
	_ = sm.TransitionTo(core.StateBuilding)

	fmt.Println("System moved to BUILDING state.")
	fmt.Println("Note: This is a skeleton implementation. Future agents (Coder Hat) will take over here.")
}
