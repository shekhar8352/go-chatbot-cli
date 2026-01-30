package llm

import "context"

// Provider defines the interface for LLM providers
type Provider interface {
	// ClassifyIntent classifies user input into one of the provided intents
	ClassifyIntent(
		ctx context.Context,
		input string,
		intents []Intent,
	) (string, error)

	// ExtractEntities extracts entities from user input based on schema
	ExtractEntities(
		ctx context.Context,
		input string,
		schema map[string]string,
	) (map[string]string, error)

	// GenerateText generates text based on a prompt
	GenerateText(
		ctx context.Context,
		prompt Prompt,
	) (string, error)
}

// Intent represents an intent with name and examples
type Intent struct {
	Name     string
	Examples []string
}

// Prompt represents a text generation prompt
type Prompt struct {
	Text string
}
