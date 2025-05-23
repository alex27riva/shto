package tui

import (
	"fmt"
	"shto/internal/types"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	hosts     []types.Host
	cursor    int
	Selected  bool
	Selection string
}

func NewModel(hosts []types.Host) Model {
	return Model{hosts: hosts}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.hosts)-1 {
				m.cursor++
			}
		case "enter":
			m.Selected = true
			m.Selection = m.hosts[m.cursor].Name
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Selected {
		return fmt.Sprintf("Connecting to %s...\n", m.Selection)
	}

	s := "Select a host to connect:\n\n"
	s += fmt.Sprintf("%-30s %-15s %-5s %-10s\n", "Host", "User", "Port", "Source")
	s += strings.Repeat("-", 65) + "\n"
	for i, host := range m.hosts {
		cursor := " " // no cursor
		line := fmt.Sprintf("%-30s %-15s %-5s %-10s", host.Name, host.Username, host.Port, host.Source)
		if m.cursor == i {
			cursor = ">" // cursor
			// Highlight the selected entry (bold and blue)
			line = fmt.Sprintf("\033[1;34m%s\033[0m", line)
		}
		s += fmt.Sprintf("%s %s\n", cursor, line)
	}
	s += "\nPress q to quit.\n"
	return s
}
