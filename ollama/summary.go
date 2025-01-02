package ollama

import "encoding/json"

type Summary struct {
	Summary string `json:"summary"`
}

func (c *Client) Summary(text string) (string, error) {

	data := ChatRequest{
		Model:    c.Model,
		Insecure: true,
		Stream:   false,
		Messages: []Message{
			{
				Role:    "system",
				Content: "Short summary of text. Return as JSON.", // llama3.2:3b
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
		},
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	cr, err := c.Chat(data)
	if err != nil {
		return "", err
	}

	var summary Summary
	err = json.Unmarshal([]byte(cr.Message.Content), &summary)
	if err != nil {
		return "", err
	}
	return summary.Summary, nil
}
