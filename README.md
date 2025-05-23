# Shto: SSH Host Selector

Shto is a command-line tool designed to simplify SSH host selection and connection. It provides a Text User Interface (TUI) for selecting SSH hosts, leveraging the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for an intuitive and interactive experience.

## Features

- **TUI for SSH Host Selection**: Easily browse and select SSH hosts using a fuzzy-search-enabled interface.
- **Integration with `.ssh/config` and `known_hosts`**: Automatically parses and merges entries from `.ssh/config` and `known_hosts`, prioritizing `.ssh/config`.
- **Username and Port Handling**: Automatically applies the correct username and port from `.ssh/config` or defaults to standard values.
- **Debugging Support**: Prints the executed SSH command for easy debugging.

## How It Works

1. **Host Parsing**:
   - Reads and parses `.ssh/config` to extract host entries, filtering out global settings (`Host *`).
   - Reads `known_hosts` to include additional hosts not listed in `.ssh/config`.
   - Merges the two sources, prioritizing `.ssh/config` entries.

2. **TUI Interaction**:
   - Displays a list of available SSH hosts with details (username, port, source).
   - Allows fuzzy searching to quickly find the desired host.

3. **SSH Connection**:
   - Executes the SSH command with the selected host's details.
   - Ensures the correct username and port are applied.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/shto.git
   cd shto
   ```

2. Build the binary:
   ```bash
   go build -o shto
   ```

3. Run the program:
   ```bash
   ./shto
   ```

## Usage

1. Run the program:
   ```bash
   ./shto
   ```

2. Use the TUI to search and select an SSH host.

3. The program will execute the SSH command and connect to the selected host.

## Example TUI Output

Below is an example of the TUI interface when running `shto`:

```
Select an SSH Host:

> host1.example.com (user: alice, port: 22, source: .ssh/config)
  host2.example.com (user: bob, port: 2222, source: known_hosts)
  host3.example.com (user: root, port: 22, source: .ssh/config)

Search: host1
```

In this example:
- `host1.example.com` is selected.
- The username is `alice` and the port is `22`, as specified in `.ssh/config`.
- The source of the host entry is displayed for clarity.

## Requirements

- Go 1.18 or later
- A valid `.ssh/config` or `known_hosts` file

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the project.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- [ssh_config](https://github.com/kevinburke/ssh_config) for parsing `.ssh/config`
- [fuzzy](https://github.com/sahilm/fuzzy) for fuzzy search functionality
