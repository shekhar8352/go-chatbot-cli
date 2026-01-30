package router

import (
	"chatbot-go/internal/bot"
	"chatbot-go/internal/llm"
	"context"
	"encoding/json"
	"fmt"
)

// LLMRouter routes using LLM for intent classification
type LLMRouter struct {
	provider llm.Provider
}

// NewLLMRouter creates a new LLM-based router
func NewLLMRouter(provider llm.Provider) *LLMRouter {
	return &LLMRouter{
		provider: provider,
	}
}

// Route uses LLM to classify intent
func (r *LLMRouter) Route(ctx context.Context, input string, intents []bot.Intent) (string, error) {
	if len(intents) == 0 {
		return "", ErrNoIntents{}
	}

	// Build intent list for LLM
	intentList := make([]llm.Intent, len(intents))
	for i, intent := range intents {
		intentList[i] = llm.Intent{
			Name:     intent.Name,
			Examples: intent.Examples,
		}
	}

	// Use LLM to classify
	intentName, err := r.provider.ClassifyIntent(ctx, input, intentList)
	if err != nil {
		return "", fmt.Errorf("LLM classification failed: %w", err)
	}

	// Validate that the returned intent exists
	for _, intent := range intents {
		if intent.Name == intentName {
			return intentName, nil
		}
	}

	return "", fmt.Errorf("LLM returned invalid intent: %s", intentName)
}

// ParseLLMResponse parses JSON response from LLM
func ParseLLMResponse(response string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, err
	}
	return result, nil
}
