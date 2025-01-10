package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type PullRequest struct {
	Model    string `json:"model"`
	Insecure bool   `json:"insecure"`
	Stream   bool   `json:"stream"`
}

type PullResponse struct {
}

func (c *Client) Pull() error {

	data := PullRequest{
		Model:    c.Model,
		Insecure: true,
		Stream:   false,
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create the request
	endpoint := c.Url + "/api/pull"
	log.Println(endpoint)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		createRecord("pull", "error")
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		createRecord("pull", "error")
		return err
	}
	defer resp.Body.Close()

	// Read the response
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// fmt.Printf("Response: %s\n", string(body))
	createRecord("pull", "success")
	return nil
}
