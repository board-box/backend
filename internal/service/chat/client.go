package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const model = "meta-llama/llama-4-scout:free"

type client struct {
	APIKey     string
	SiteURL    string
	SiteTitle  string
	HTTPClient *http.Client
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type response struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func newClient(apiKey string) *client {
	return &client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *client) chat(messages []message) (message, error) {
	reqBody := request{
		Model:    model,
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return message{}, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return message{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	if c.SiteURL != "" {
		req.Header.Set("HTTP-Referer", c.SiteURL)
	}
	if c.SiteTitle != "" {
		req.Header.Set("X-Title", c.SiteTitle)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return message{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return message{}, err
	}

	if resp.StatusCode != 200 {
		return message{}, fmt.Errorf("API error: %s", string(respBytes))
	}

	var result response
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return message{}, err
	}

	return result.Choices[0].Message, nil
}
