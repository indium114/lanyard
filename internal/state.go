package internal

import (
	"log"
	"os"
	"strings"
)

func StateDir() string {
	home, err := os.UserConfigDir()
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
	if err := os.WriteFile(StateDir()+"/state", bytes, 0o700); err != nil {
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
