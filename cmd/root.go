package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"shto/internal/ssh"
	"shto/internal/tui"
	"shto/internal/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

var sshUser string

// NewRootCmd creates the root Cobra command for the application
var rootCmd = &cobra.Command{
	Use:   "shto",
	Short: "SSH host selector",
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		if sshUser == "" {
			sshUser = os.Getenv("SHTO_SSH_USER")
		}
		if sshUser == "" {
			sshUser = usr.Username
		}

		if sshUser != "" {
			fmt.Printf("Using SSH user: %s\n", sshUser)
		} else {
			fmt.Println("No SSH user specified. Falling back to default.")
		}

		knownHostsPath := filepath.Join(usr.HomeDir, ".ssh", "known_hosts")
		file, err := os.Open(knownHostsPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		hostsMap := make(map[string]struct{})
		hostPorts := make(map[string]string)
		scanner := bufio.NewScanner(file)

		// Updated parsing logic to handle ports in known_hosts entries
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "|") || line == "" {
				continue
			}
			parts := strings.Split(line, " ")
			hostnames := strings.Split(parts[0], ",")
			for _, host := range hostnames {
				var hostname, port string
				if strings.HasPrefix(host, "[") && strings.Contains(host, "]:") {
					// Extract hostname and port from [host]:port format
					endIdx := strings.Index(host, "]")
					hostname = host[1:endIdx]
					port = host[endIdx+2:]
				} else {
					hostname = host
					port = "22" // Default port
				}
				hostsMap[hostname] = struct{}{}
				if port != "" {
					// Store port information in a separate map or structure if needed
					hostPorts[hostname] = port
				}
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		var hosts []string
		for h := range hostsMap {
			hosts = append(hosts, h)
		}
		sort.Strings(hosts)

		if len(hosts) == 0 {
			fmt.Println("No hosts found in known_hosts.")
			return
		}

		var input string
		fmt.Print("Enter search term (leave empty to show all): ")
		fmt.Scanln(&input)
		if input != "" {
			hosts = filterHosts(input, hosts)
		}

		sshConfig := tui.ParseSSHConfig()
		for host := range sshConfig {
			if host == "*" {
				delete(sshConfig, host)
			}
		}

		var hostEntries []types.Host
		for host, config := range sshConfig {
			username := sshUser
			if user, exists := config["User"]; exists && user != "" {
				username = user
			}
			port := config["Port"]
			if port == "" {
				port = "22"
			}
			hostEntries = append(hostEntries, types.Host{
				Name:     host,
				IP:       host,
				Username: username,
				Port:     port,
				Source:   "config",
			})
		}

		knownHostsSet := make(map[string]struct{})
		for _, entry := range hostEntries {
			knownHostsSet[entry.Name] = struct{}{}
		}

		// Update hostEntries to include the parsed port from known_hosts
		for _, h := range hosts {
			if _, exists := knownHostsSet[h]; !exists {
				port := hostPorts[h] // Retrieve the port from the parsed hostPorts map
				if port == "" {
					port = "22" // Default port
				}
				hostEntries = append(hostEntries, types.Host{
					Name:     h,
					IP:       h,
					Username: sshUser,
					Port:     port,
					Source:   "known_hosts",
				})
			}
		}

		p := tea.NewProgram(tui.NewModel(hostEntries))
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			return
		}

		finalModel := m.(tui.Model)
		if finalModel.Selected {
			selectedHost := finalModel.Selection
			var username, port string

			if config, ok := sshConfig[selectedHost]; ok {
				if user, exists := config["User"]; exists {
					username = user
				}
				if p, exists := config["Port"]; exists {
					port = p
				}
			}

			if username == "" {
				username = sshUser
			}
			if port == "" {
				port = "22"
			}

			ssh.Connect(username, selectedHost, port)
		}
	},
}

func filterHosts(input string, hosts []string) []string {
	matches := fuzzy.Find(input, hosts)
	filtered := make([]string, len(matches))
	for i, match := range matches {
		filtered[i] = hosts[match.Index]
	}
	return filtered
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return nil
}

func init() {
	rootCmd.Flags().StringVarP(&sshUser, "user", "u", "", "Specify the SSH user")
}
