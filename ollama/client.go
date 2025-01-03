package ollama

type Client struct {
	Url   string
	Model string
}

func NewClient(url, model string) *Client {
	return &Client{
		Url:   url,
		Model: model,
	}
}

func NewClientFromConfig() *Client {
	return &Client{
		Url:   "http://localhost:11434",
		Model: "llama3.2:3b",
		// Model: "llama3.2:1b",
	}
}
