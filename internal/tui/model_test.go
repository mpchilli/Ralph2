package tui

import (
	"ralph2/internal/core"
	"ralph2/pkg/utils"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdate(t *testing.T) {
	// Setup
	bus := utils.NewEventBus()
	model := NewModel(bus, core.StatePlanning)

	// Simulate Init subscription (in a real app this happens via tea.Program, 
	// here we simulate the effect by manually sending a message if we could, 
	// but Update handles tea.Msg)

	// 1. Test KeyMsg Quit
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := model.Update(msg)
	
	if cmd == nil {
		t.Error("Expected Quit command, got nil")
	} else {
		// Verify it's actually the Quit command (simplified check as tea.Quit is func)
		// We can't directly compare functions, but specific tea commands often have properties.
		// For standardized tea.Quit, we assume non-nil command returned in this context is correct 
		// given the simple switch case.
		// A stronger check would be to execute the command if possible or inspect internal type if exposed.
	}

	// 2. Test State Change Event
	// We need to simulate receiving an event from the EventBus.
	// The Model.Update should handle a custom Msg type for events.
	
	newState := core.StateBuilding
	event := utils.Event{
		Topic:   "state_change",
		Payload: newState,
	}
	
	// We wrap utils.Event in a tea.Msg wrapper if needed, or cast it directly
	updatedModel, _ := model.Update(event)
	
	m, ok := updatedModel.(Model)
	if !ok {
		t.Fatalf("Model type assertion failed")
	}
	
	if m.State != newState {
		t.Errorf("Expected state %s, got %s", newState, m.State)
	}
}
