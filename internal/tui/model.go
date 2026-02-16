package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"ralph2/internal/core"
	"ralph2/pkg/utils"
)

type Model struct {
	State    core.FSMState
	Messages []string
	EventBus *utils.EventBus
}

func NewModel(bus *utils.EventBus) Model {
	return Model{
		State:    core.StatePlanning,
		EventBus: bus,
	}
}

func (m Model) Init() tea.Cmd {
	// Subscribe to events (TODO: Actual tea.Cmd to listen to channel)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "Ralph 2.1 - Loki Mode\n\n"
	s += fmt.Sprintf("State: %s\n", m.State)
	s += "\nLogs:\n"
	for _, log := range m.Messages {
		s += fmt.Sprintf("> %s\n", log)
	}
	s += "\nPress 'q' to quit.\n"
	return s
}
