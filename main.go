package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	hosts     []string
	cursor    int
	selected  bool
	selection string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.selected = true
			m.selection = m.hosts[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.selected {
		return fmt.Sprintf("Connecting to %s...\n", m.selection)
	}

	s := "Select a host to connect:\n\n"
	for i, host := range m.hosts {
		cursor := " " // no cursor
		line := host  // default line text
		if m.cursor == i {
			cursor = ">" // cursor
			// Highlight the selected entry (bold and blue)
			line = fmt.Sprintf("\033[1;34m%s\033[0m", host)
		}
		s += fmt.Sprintf("%s %s\n", cursor, line)
	}
	s += "\nPress q to quit.\n"
	return s
}

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	knownHostsPath := filepath.Join(usr.HomeDir, ".ssh", "known_hosts")
	file, err := os.Open(knownHostsPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hostsMap := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "|") || line == "" {
			continue // Skip hashed or empty lines
		}
		parts := strings.Split(line, " ")
		hostnames := strings.Split(parts[0], ",")
		for _, host := range hostnames {
			hostsMap[host] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var hosts []string
	for h := range hostsMap {
		hosts = append(hosts, h)
	}

	if len(hosts) == 0 {
		fmt.Println("No hosts found in known_hosts.")
		return
	}

	// Initialize the Bubble Tea program
	p := tea.NewProgram(model{hosts: hosts})

	// Run the TUI
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}

	// Handle the selected host
	finalModel := m.(model)
	if finalModel.selected {
		cmd := exec.Command("ssh", finalModel.selection)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("SSH failed: %v\n", err)
		}
	}
}
