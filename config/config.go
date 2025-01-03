package config

import (
	"fmt"
	"os"
	"strings"
)

// returns value of JOURNAL_OLLAMA_CONFIG or journal_config
func ConfigDir() string {
	val, ok := os.LookupEnv("JOURNAL_OLLAMA_CONFIG")
	if !ok || val == "" {
		return "journal_config"
	}
	return val
}

// returns value of JOURNAL_OLLAMA_DATA or journal_data
func DataDir() string {
	val, ok := os.LookupEnv("JOURNAL_OLLAMA_DATA")
	if !ok || val == "" {
		return "journal_data"
	}
	return val
}

// returns true iff JOURNAL_OLLAMA_INSECURE is truthy
func OllamaInsecure() bool {
	val, ok := os.LookupEnv("JOURNAL_OLLAMA_INSECURE")
	if ok {
		lower := strings.ToLower(val)
		if lower == "1" || lower == "on" || lower == "yes" || lower == "true" {
			return true
		}
	}
	return false
}

// returns JOURNAL_OLLAMA_URL or panics
func OllamaUrl() string {
	key := "JOURNAL_OLLAMA_URL"
	val, ok := os.LookupEnv(key)
	if ok && val != "" {
		return val
	} else {
		panic(fmt.Sprintf("Set %s", key))
	}
}

// returns JOURNAL_PASSWORD or panics
func Password() string {
	key := "JOURNAL_PASSWORD"
	val, ok := os.LookupEnv(key)
	if ok && val != "" {
		return val
	} else {
		panic(fmt.Sprintf("Set %s", key))
	}
}
