package config

import (
	"fmt"
	"os"
	"strings"
)

// returns value of JOURNAL_CONFIG_DIR or journal_config
func ConfigDir() string {
	val, ok := os.LookupEnv("JOURNAL_CONFIG_DIR")
	if !ok || val == "" {
		return "journal_config"
	}
	return val
}

// returns value of JOURNAL_DATA_DIR or journal_data
func DataDir() string {
	val, ok := os.LookupEnv("JOURNAL_DATA_DIR")
	if !ok || val == "" {
		return "journal_data"
	}
	return val
}

// returns value of JOURNAL_PORT or 8080
func Port() string {
	val, ok := os.LookupEnv("JOURNAL_PORT")
	if !ok || val == "" {
		return "8080"
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

// returns JOURNAL_SESSION_KEY or panics
func SessionKey() string {
	key := "JOURNAL_SESSION_KEY"
	val, ok := os.LookupEnv(key)
	if ok && val != "" {
		return val
	} else {
		panic(fmt.Sprintf("Set %s", key))
	}
}

// returns true iff JOURNAL_SESSION_SECURE is truthy
func SessionSecure() bool {
	val, ok := os.LookupEnv("JOURNAL_SESSION_SECURE")
	if ok {
		lower := strings.ToLower(val)
		if lower == "0" || lower == "off" || lower == "no" || lower == "false" {
			return false
		}
	}
	return true
}
