package ollama

import (
	"encoding/json"
)

type Keywords struct {
	Keywords []string `json:"keywords"`
}

func (c *Client) Keywords(text string) ([]string, error) {

	data := ChatRequest{
		Model:    c.Model,
		Insecure: true,
		Stream:   false,
		Messages: []Message{
			{
				Role: "system",
				// Content: "You produce keywords for text. Return as JSON.",
				Content: "Produce the most important keywords for text. Return as JSON.", // llama3.2:3b
			},
			{
				Role:    "user",
				Content: text,
			},
		},
		Format: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"keywords": map[string]interface{}{
					"type": "array",
					"items": map[string]string{
						"type": "string",
					},
				},
			},
			"required": []string{
				"keywords",
			},
		},
		// Format: map[string]interface{}{
		// 	"type": "array",
		// 	"items": map[string]interface{}{
		// 		"type": "string",
		// 	},
		// },
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	cr, err := c.Chat(data)
	if err != nil {
		return nil, err
	}

	var keywords Keywords
	// keywords := []string{}
	err = json.Unmarshal([]byte(cr.Message.Content), &keywords)
	if err != nil {
		return nil, err
	}
	return keywords.Keywords, nil
	// return keywords, nil
}
