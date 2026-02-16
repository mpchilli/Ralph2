package core

import (
	"fmt"
	"sync"
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

// EventType for the EventBus
type EventType string

const (
	EventStateChange EventType = "state_change"
	EventLog         EventType = "log"
)

// Event represents a system event
type Event struct {
	Type    EventType
	Payload interface{}
}

// StateManager handles the FSM
type StateManager struct {
	CurrentState FSMState
	mu           sync.RWMutex
	eventChan    chan Event
}

// NewStateManager creates a new FSM
func NewStateManager() *StateManager {
	return &StateManager{
		CurrentState: StatePlanning,
		eventChan:    make(chan Event, 100),
	}
}

// TransitionTo changes state safely
func (sm *StateManager) TransitionTo(newState FSMState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// TODO: Add validity checks (e.g., cannot go from FAILED to COMPLETE directly)
	fmt.Printf("Transitioning from %s to %s\n", sm.CurrentState, newState)
	
	sm.CurrentState = newState
	
	// Broadcast
	sm.eventChan <- Event{
		Type:    EventStateChange,
		Payload: newState,
	}
	
	return nil
}

// GetState returns current state
func (sm *StateManager) GetState() FSMState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.CurrentState
}
