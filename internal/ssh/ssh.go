package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Connect(user, host, port string) {
	args := []string{fmt.Sprintf("%s@%s", user, host)}
	if port != "" {
		args = append(args, "-p", port)
	}

	// Print the SSH command for debugging
	fmt.Printf("Executing SSH command: ssh %s\n", strings.Join(args, " "))

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("SSH failed: %v\n", err)
	}
}
