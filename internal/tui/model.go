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
	sub      <-chan utils.Event
}

func NewModel(bus *utils.EventBus, initialState core.FSMState) Model {
	return Model{
		State:    initialState,
		EventBus: bus,
		sub:      bus.Subscribe("state_change"),
	}
}

func (m Model) Init() tea.Cmd {
	return waitForEvent(m.sub)
}

func waitForEvent(sub <-chan utils.Event) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Handle window resize if needed (e.g., update layout dimensions)
		return m, nil
	case utils.Event:
		if msg.Topic == "state_change" {
			if newState, ok := msg.Payload.(core.FSMState); ok {
				m.State = newState
			}
		}
		return m, waitForEvent(m.sub)
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
