package internal

import (
	"log"
	"os"
	"os/exec"
)

func Wipe(sock, pid string) {
	// remove from agent
	cmd := exec.Command("ssh-add", "-D")
	cmd.Stdout = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, "SSH_AUTH_SOCK="+sock, "SSH_AGENT_PID="+pid)

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to remove keys from ssh-agent")
	}

	WriteUnlocked([]string{})
}
