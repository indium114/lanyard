package internal

import (
	"log"
	"os"
	"strings"
)

func StateDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return home + "/.lanyard"
}

func InitStateDir() {
	if err := os.MkdirAll(StateDir(), 0o700); err != nil {
		log.Fatal(err)
	}
}

func WriteStateFile(sock, pid string) {
	bytes := []byte(sock + " " + pid)
	if err := os.WriteFile(StateDir()+"/state", bytes, 0o600); err != nil {
		log.Fatal(err)
	}
}

func ReadStateFile() (bool, string, string) {
	content, err := os.ReadFile(StateDir() + "/state")
	if err != nil {
		return false, "", ""
	}

	state := strings.Split(string(content), " ")

	return true, state[0], state[1]
}

// MARK: unlocked state
func WriteUnlocked(keys []string) {
	path := StateDir() + "/unlocked"
	content := strings.Join(keys, "\n")

	_ = os.WriteFile(path, []byte(content), 0o600)
}

func ReadUnlocked() []string {
	path := StateDir() + "/unlocked"
	content, err := os.ReadFile(path)
	if err != nil {
		return []string{}
	}

	text := strings.TrimRight(string(content), "\n")
	if text == "" {
		return []string{}
	}

	return strings.Split(text, "\n")
}
