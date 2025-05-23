package tui

import (
	"os"
	"path/filepath"

	"github.com/kevinburke/ssh_config"
)

// Exported the `parseSSHConfig` function for use in other packages
func ParseSSHConfig() map[string]map[string]string {
	configPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil // Return an empty map if the file doesn't exist
	}
	defer file.Close()

	parsedConfig := make(map[string]map[string]string)

	sshConfig, err := ssh_config.Decode(file)
	if err != nil {
		return nil // Return an empty map if decoding fails
	}

	for _, host := range sshConfig.Hosts {
		if len(host.Patterns) > 0 {
			hostName := host.Patterns[0].String()
			parsedConfig[hostName] = make(map[string]string)
			for _, node := range host.Nodes {
				if kv, ok := node.(*ssh_config.KV); ok {
					parsedConfig[hostName][kv.Key] = kv.Value
				}
			}
		}
	}

	return parsedConfig
}
