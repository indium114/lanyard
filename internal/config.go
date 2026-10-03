package internal

import (
	"log"
	"os"
	"strings"
)

func configPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to load config. Does it exist?", err)
	}

	return configDir + "/lanyard/keys.conf"
}

func LoadConfig() []string {
	content, err := os.ReadFile(configPath())
	if err != nil {
		log.Fatal(err)
	}

	text := strings.TrimRight(string(content), "\n")
	if text == "" {
		return []string{}
	}

	return strings.Split(text, "\n")
}
