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

// StateManager handles the FSM
type StateManager struct {
	CurrentState FSMState
	mu           sync.RWMutex
	bus          *utils.EventBus
}

// NewStateManager creates a new FSM
func NewStateManager(bus *utils.EventBus) *StateManager {
	return &StateManager{
		CurrentState: StatePlanning,
		bus:          bus,
	}
}

// TransitionTo changes state safely
func (sm *StateManager) TransitionTo(newState FSMState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// TODO: Add validity checks (e.g., cannot go from FAILED to COMPLETE directly)
	fmt.Printf("Transitioning from %s to %s\n", sm.CurrentState, newState)
	
	sm.CurrentState = newState
	
	// Broadcast via EventBus
	if sm.bus != nil {
		sm.bus.Publish("state_change", newState)
	}
	
	return nil
}

// GetState returns current state
func (sm *StateManager) GetState() FSMState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.CurrentState
}
