package ollama

import (
	"github.com/cwpearson/journal/config"
)

type Client struct {
	Url      string
	Model    string
	Insecure bool
}

func NewClientFromConfig() *Client {
	return &Client{
		Url:   config.OllamaUrl(),
		Model: "llama3.2:3b",
		// Model: "llama3.2:1b",
		Insecure: config.OllamaInsecure(),
	}
}
