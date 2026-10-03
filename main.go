package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"
	"github.com/indium114/lanyard/internal"
)

func startAgent() (string, string) {
	_ = os.RemoveAll(internal.StateDir())
	internal.InitStateDir()
	cmd := exec.Command("ssh-agent", "-a", internal.StateDir()+"/agent.sock")
	outBytes, err := cmd.Output()
	out := string(outBytes)
	lines := strings.Split(out, "\n")
	if err != nil {
		log.Fatal("Failed to start ssh-agent due to ", err)
	}

	sock := strings.TrimSuffix(strings.TrimPrefix(lines[0], "SSH_AUTH_SOCK="), "; export SSH_AUTH_SOCK;")
	pid := strings.TrimSuffix(strings.TrimPrefix(lines[1], "SSH_AGENT_PID="), "; export SSH_AGENT_PID;")

	internal.WriteStateFile(sock, pid)
	return sock, pid
}

func bootstrap() (string, string) {
	internal.InitStateDir()

	// decide whether or not we need to start a new agent
	success, _, pidString := internal.ReadStateFile()
	pid, _ := strconv.Atoi(pidString)

	exists := false
	proc, err := os.FindProcess(pid)
	if err != nil {
		exists = false
	}
	if err := proc.Signal(syscall.Signal(0)); err == nil {
		exists = true
	}

	if !success || !exists {
		startAgent()
	}

	_, sock, newPid := internal.ReadStateFile()
	return sock, newPid
}

func main() {
	sock, pid := bootstrap()
	path := internal.StateDir() + "/activate.nu"

	nuScript := fmt.Sprintf("{SSH_AGENT_PID: \"%s\", SSH_AUTH_SOCK: \"%s\"}", pid, sock)
	os.WriteFile(path, []byte(nuScript), 0o700)

	// confirm whether or not to add keys
	if len(internal.LoadConfig()) == len(internal.ReadUnlocked()) {
		os.Exit(0)
	}
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Unlock Keys?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if confirm {
		home, _ := os.UserHomeDir()

		var keys []string
		for _, configuredKey := range internal.LoadConfig() {
			found := false
			for _, unlockedKey := range internal.ReadUnlocked() {
				if configuredKey == unlockedKey {
					found = true
					break
				}
			}
			if !found {
				keys = append(keys, configuredKey)
			}
		}

		for _, key := range keys {
			cmd := exec.Command("ssh-add", home+"/.ssh/"+key)
			cmd.Stdout = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Env = append(cmd.Env, "SSH_AUTH_SOCK="+sock, "SSH_AGENT_PID="+pid)

			if err := cmd.Run(); err == nil {
				internal.WriteUnlocked(append(internal.ReadUnlocked(), key))
			}
		}
	}
}
