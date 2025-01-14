package ollama

import "encoding/json"

type Summary struct {
	Summary string `json:"summary"`
}

func (c *Client) Summary(text string) (string, error) {

	if err := c.Pull(); err != nil {
		return "", err
	}

	data := ChatRequest{
		Model:    c.Model,
		Insecure: true,
		Stream:   false,
		Messages: []Message{
			{
				Role: "system",
				// Content: "Short summary of text. Return as JSON.", // llama3.2:3b
				Content: "Very short summary of diary entry, for the author. Address author as \"you.\" Return as JSON.", // llama3.2:3b
			},
			{
				Role:    "user",
				Content: text,
			},
		},
		Format: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"summary": map[string]string{
					"type": "string",
				},
			},
			"required": []string{
				"summary",
			},
		},
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	cr, err := c.Chat(data)
	if err != nil {
		createRecord("summary", "error")
		return "", err
	}

	var summary Summary
	err = json.Unmarshal([]byte(cr.Message.Content), &summary)
	if err != nil {
		createRecord("summary", "error")
		return "", err
	}
	createRecord("summary", "success")
	return summary.Summary, nil
}
