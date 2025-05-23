package types

// Host represents an SSH host with its details
type Host struct {
	Name     string
	IP       string
	Username string
	Port     string
	Source   string
}
