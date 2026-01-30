package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaProvider is an HTTP-based stub for Ollama integration
type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama2"
	}

	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ClassifyIntent uses Ollama to classify intent
func (o *OllamaProvider) ClassifyIntent(
	ctx context.Context,
	input string,
	intents []Intent,
) (string, error) {
	// Build prompt
	prompt := o.buildIntentClassificationPrompt(input, intents)

	// Call Ollama API
	response, err := o.callAPI(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("ollama API call failed: %w", err)
	}

	// Parse JSON response
	var result struct {
		Intent string `json:"intent"`
	}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return "", fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return result.Intent, nil
}

// ExtractEntities uses Ollama to extract entities
func (o *OllamaProvider) ExtractEntities(
	ctx context.Context,
	input string,
	schema map[string]string,
) (map[string]string, error) {
	prompt := o.buildEntityExtractionPrompt(input, schema)
	response, err := o.callAPI(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("ollama API call failed: %w", err)
	}

	var result map[string]string
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return result, nil
}

// GenerateText uses Ollama to generate text
func (o *OllamaProvider) GenerateText(
	ctx context.Context,
	prompt Prompt,
) (string, error) {
	return o.callAPI(ctx, prompt.Text)
}

// buildIntentClassificationPrompt builds a prompt for intent classification
func (o *OllamaProvider) buildIntentClassificationPrompt(input string, intents []Intent) string {
	var buf bytes.Buffer
	buf.WriteString("Classify the following user input into one of the provided intents.\n")
	buf.WriteString("Respond with JSON only: {\"intent\": \"<intent_name>\"}\n\n")
	buf.WriteString("User input: " + input + "\n\n")
	buf.WriteString("Available intents:\n")
	for _, intent := range intents {
		buf.WriteString(fmt.Sprintf("- %s (examples: %v)\n", intent.Name, intent.Examples))
	}
	return buf.String()
}

// buildEntityExtractionPrompt builds a prompt for entity extraction
func (o *OllamaProvider) buildEntityExtractionPrompt(input string, schema map[string]string) string {
	var buf bytes.Buffer
	buf.WriteString("Extract entities from the following user input.\n")
	buf.WriteString("Respond with JSON only containing the extracted entities.\n\n")
	buf.WriteString("User input: " + input + "\n\n")
	buf.WriteString("Schema:\n")
	for key, desc := range schema {
		buf.WriteString(fmt.Sprintf("- %s: %s\n", key, desc))
	}
	return buf.String()
}

// callAPI makes an HTTP request to Ollama API
func (o *OllamaProvider) callAPI(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", o.baseURL)

	payload := map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}
