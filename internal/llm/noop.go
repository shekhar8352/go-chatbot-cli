package llm

import (
	"context"
	"errors"
)

// NoopProvider is a no-op LLM provider that always returns errors
// This ensures the system works without any LLM configured
type NoopProvider struct{}

// NewNoopProvider creates a new no-op provider
func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

// ClassifyIntent always returns an error (no LLM available)
func (n *NoopProvider) ClassifyIntent(
	ctx context.Context,
	input string,
	intents []Intent,
) (string, error) {
	return "", errors.New("no LLM provider configured")
}

// ExtractEntities always returns an error (no LLM available)
func (n *NoopProvider) ExtractEntities(
	ctx context.Context,
	input string,
	schema map[string]string,
) (map[string]string, error) {
	return nil, errors.New("no LLM provider configured")
}

// GenerateText always returns an error (no LLM available)
func (n *NoopProvider) GenerateText(
	ctx context.Context,
	prompt Prompt,
) (string, error) {
	return "", errors.New("no LLM provider configured")
}
