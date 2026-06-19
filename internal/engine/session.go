package engine

import (
	"fmt"
	"strconv"
	"strings"

	"chatbot-go/internal/bot"
)

// SetVariable sets a variable in the session
func (e *Engine) SetVariable(key, value string) {
	if e.session.Variables == nil {
		e.session.Variables = make(map[string]string)
	}
	e.session.Variables[key] = value
}

// GetVariable retrieves a variable from the session
func (e *Engine) GetVariable(key string) (string, bool) {
	if e.session.Variables == nil {
		return "", false
	}
	val, exists := e.session.Variables[key]
	return val, exists
}

// AppendVariable appends a value to an existing variable.
func (e *Engine) AppendVariable(key, value, separator string) {
	if e.session.Variables == nil {
		e.session.Variables = make(map[string]string)
	}

	existing, ok := e.session.Variables[key]
	if !ok || existing == "" {
		e.session.Variables[key] = value
		return
	}

	e.session.Variables[key] = existing + separator + value
}

// IncrementVariable increments a numeric session variable.
func (e *Engine) IncrementVariable(key string, by int) error {
	if e.session.Variables == nil {
		e.session.Variables = make(map[string]string)
	}

	current := e.session.Variables[key]
	if current == "" {
		current = "0"
	}

	n, err := strconv.Atoi(strings.TrimSpace(current))
	if err != nil {
		return fmt.Errorf("variable '%s' is not numeric: %w", key, err)
	}

	e.session.Variables[key] = strconv.Itoa(n + by)
	return nil
}

// DeleteVariable removes a variable from the session.
func (e *Engine) DeleteVariable(key string) {
	if e.session.Variables == nil {
		return
	}
	delete(e.session.Variables, key)
}

// AddTurn adds a turn to the conversation history
func (e *Engine) AddTurn(node, userInput, response string) {
	e.session.History = append(e.session.History, Turn{
		Node:      node,
		UserInput: userInput,
		Response:  response,
	})
}

// EffectiveIntents returns node intents followed by global intents.
func (e *Engine) EffectiveIntents(node *bot.Node) []bot.Intent {
	intents := make([]bot.Intent, 0, len(node.Intents)+len(e.bot.GlobalIntents))
	intents = append(intents, node.Intents...)
	intents = append(intents, e.bot.GlobalIntents...)
	return intents
}
