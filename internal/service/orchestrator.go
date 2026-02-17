package service

import (
	"fmt"
	"ralph2/internal/coder"
	"ralph2/internal/core"
	"ralph2/internal/git"
	"ralph2/internal/planner"
	"ralph2/pkg/utils"
)

type OrchestratorService struct {
	Bus *utils.EventBus
	SM  *core.StateManager
}

func NewOrchestratorService() *OrchestratorService {
	bus := utils.NewEventBus()
	sm := core.NewStateManager(bus)
	return &OrchestratorService{
		Bus: bus,
		SM:  sm,
	}
}

func (s *OrchestratorService) Run(prompt string, complexity string) error {
	fmt.Printf("🚀 Starting Ralph2 with prompt: %q\n", prompt)
	fmt.Printf("Complexity: %s\n", complexity)

	// Step 0: Git Isolation Check
	isDirty, err := git.CheckDirty()
	if err != nil {
		return fmt.Errorf("error checking git status: %v", err)
	}
	if isDirty {
		return fmt.Errorf("working tree is dirty. Please commit or stash changes before running")
	}

	branch, err := git.CreateTaskBranch()
	if err != nil {
		return fmt.Errorf("error creating task branch: %v", err)
	}
	fmt.Printf("🌿 Switched to task branch: %s\n", branch)

	// Validate and set complexity
	level := core.ComplexityLevel(complexity)
	switch level {
	case core.ComplexityFast, core.ComplexityStreamlined, core.ComplexityFull:
		s.SM.SetComplexity(level)
	default:
		fmt.Printf("Warning: Unknown complexity %q, defaulting to streamlined\n", complexity)
	}

	// Move to planning
	s.SM.TransitionTo(core.StatePlanning)

	// Execute Planning
	pln := planner.NewHeuristicPlanner("spec.md")
	if err := pln.Plan(prompt); err != nil {
		_ = s.SM.TransitionTo(core.StateFailed)
		return fmt.Errorf("error during planning: %v", err)
	}

	fmt.Println("Plan generated successfully in spec.md")

	// Commit Plan
	if err := git.CommitChanges(fmt.Sprintf("docs: generated plan for %q", prompt)); err != nil {
		fmt.Printf("Warning: Failed to commit plan: %v\n", err)
	}

	// Move to building
	_ = s.SM.TransitionTo(core.StateBuilding)

	fmt.Println("System moved to BUILDING state.")

	// Execute Building
	cdr := coder.NewMockCoder("hello.go")
	if err := cdr.Build("spec.md"); err != nil {
		_ = s.SM.TransitionTo(core.StateFailed)
		return fmt.Errorf("error during building: %v", err)
	}

	fmt.Println("Code generated successfully in hello.go")

	// Commit Code
	if err := git.CommitChanges(fmt.Sprintf("feat: implement logic for %q", prompt)); err != nil {
		fmt.Printf("Warning: Failed to commit code: %v\n", err)
	}

	// Move to verifying
	_ = s.SM.TransitionTo(core.StateVerifying)

	fmt.Println("System moved to VERIFYING state.")

	// Simulated loop check
	lg := core.NewLoopGuardian()
	if isLoop, _ := lg.RecordFailure("simulated error"); isLoop {
		fmt.Println("🛑 Loop detection triggered!")
	} else {
		fmt.Println("✅ No infinite loops detected.")
	}

	fmt.Println("Note: This is a skeleton implementation. Future agents (Tester Hat) will take over here.")
	return nil
}
