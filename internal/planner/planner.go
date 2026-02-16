package planner

import (
	"fmt"
	"os"
)

// Planner defines the interface for generating a specification
type Planner interface {
	Plan(prompt string) error
}

// HeuristicPlanner is a simple implementation that mocks LLM planning
type HeuristicPlanner struct {
	OutputFile string
}

// NewHeuristicPlanner creates a new heuristic planner
func NewHeuristicPlanner(outputFile string) *HeuristicPlanner {
	return &HeuristicPlanner{
		OutputFile: outputFile,
	}
}

// Plan generates a spec file based on the prompt
func (hp *HeuristicPlanner) Plan(prompt string) error {
	content := fmt.Sprintf("# Specification\n\nGenerated from prompt: %s\n\n## Tasks\n- [ ] Implement core logic\n- [ ] Verify functionality\n", prompt)
	
	err := os.WriteFile(hp.OutputFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write spec: %w", err)
	}
	
	return nil
}
