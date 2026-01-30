package actions

import (
	"chatbot-go/internal/bot"
	"fmt"
)

// SessionMutator defines the interface for mutating session state
type SessionMutator interface {
	SetVariable(key, value string)
}

// Executor executes actions on the session
type Executor struct {
	mutator SessionMutator
}

// NewExecutor creates a new action executor
func NewExecutor(mutator SessionMutator) *Executor {
	return &Executor{
		mutator: mutator,
	}
}

// Execute executes a single action
func (ex *Executor) Execute(action bot.Action, userInput string) error {
	switch action.Type {
	case "set_var":
		return ex.executeSetVar(action, userInput)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// executeSetVar executes a set_var action
func (ex *Executor) executeSetVar(action bot.Action, userInput string) error {
	// If args has a "value" key, use it; otherwise use userInput
	value := userInput
	if val, ok := action.Args["value"]; ok {
		if strVal, ok := val.(string); ok {
			value = strVal
		} else {
			return fmt.Errorf("set_var value must be a string")
		}
	}

	// Get the variable name from args
	varName, ok := action.Args["name"].(string)
	if !ok {
		// Fallback: check if there's a "save_as" pattern
		// This handles the case where input.save_as is used
		return fmt.Errorf("set_var requires 'name' argument")
	}

	ex.mutator.SetVariable(varName, value)
	return nil
}
