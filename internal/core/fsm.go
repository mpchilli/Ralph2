package core

import (
	"fmt"
	"sync"
	"ralph2/pkg/utils"
)

// FSMState represents the state of the orchestrator
type FSMState string

const (
	StatePlanning   FSMState = "PLANNING"
	StateBuilding   FSMState = "BUILDING"
	StateVerifying  FSMState = "VERIFYING"
	StateReview     FSMState = "REVIEW"
	StateComplete   FSMState = "COMPLETE"
	StateFailed     FSMState = "FAILED"
)

// ComplexityLevel represents the depth of the loop
type ComplexityLevel string

const (
	ComplexityFast        ComplexityLevel = "fast"
	ComplexityStreamlined ComplexityLevel = "streamlined"
	ComplexityFull        ComplexityLevel = "full"
)

// StateManager handles the FSM
type StateManager struct {
	CurrentState FSMState
	Complexity   ComplexityLevel
	mu           sync.RWMutex
	bus          *utils.EventBus
}

// NewStateManager creates a new StateManager
func NewStateManager(bus *utils.EventBus) *StateManager {
	return &StateManager{
		CurrentState: StatePlanning,
		Complexity:   ComplexityStreamlined,
		bus:          bus,
	}
}

// TransitionTo changes state safely
func (sm *StateManager) TransitionTo(newState FSMState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// TODO: Add validity checks
	fmt.Printf("Transitioning from %s to %s\n", sm.CurrentState, newState)
	
	sm.CurrentState = newState
	
	// Broadcast via EventBus
	if sm.bus != nil {
		sm.bus.Publish("state_change", string(newState))
	}
	
	return nil
}

// SetComplexity updates the complexity level
func (sm *StateManager) SetComplexity(level ComplexityLevel) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.Complexity = level
}

// GetState returns current state
func (sm *StateManager) GetState() FSMState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.CurrentState
}

// GetComplexity returns current complexity
func (sm *StateManager) GetComplexity() ComplexityLevel {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.Complexity
}
