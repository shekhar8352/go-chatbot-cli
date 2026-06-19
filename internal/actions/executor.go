package actions

import (
	"chatbot-go/internal/bot"
	"fmt"
)

// SessionMutator defines the interface for mutating session state
type SessionMutator interface {
	SetVariable(key, value string)
	GetVariable(key string) (string, bool)
	AppendVariable(key, value, separator string)
	IncrementVariable(key string, by int) error
	DeleteVariable(key string)
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
	case "append_var":
		return ex.executeAppendVar(action, userInput)
	case "copy_var":
		return ex.executeCopyVar(action)
	case "increment_var":
		return ex.executeIncrementVar(action)
	case "delete_var":
		return ex.executeDeleteVar(action)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (ex *Executor) executeSetVar(action bot.Action, userInput string) error {
	value := userInput
	if val, ok := action.Args["value"]; ok {
		strVal, ok := val.(string)
		if !ok {
			return fmt.Errorf("set_var value must be a string")
		}
		value = strVal
	}

	varName, ok := action.Args["name"].(string)
	if !ok {
		return fmt.Errorf("set_var requires 'name' argument")
	}

	ex.mutator.SetVariable(varName, value)
	return nil
}

func (ex *Executor) executeAppendVar(action bot.Action, userInput string) error {
	varName, ok := action.Args["name"].(string)
	if !ok {
		return fmt.Errorf("append_var requires 'name' argument")
	}

	value := userInput
	if val, ok := action.Args["value"]; ok {
		strVal, ok := val.(string)
		if !ok {
			return fmt.Errorf("append_var value must be a string")
		}
		value = strVal
	}

	separator := ", "
	if sep, ok := action.Args["separator"].(string); ok {
		separator = sep
	}

	ex.mutator.AppendVariable(varName, value, separator)
	return nil
}

func (ex *Executor) executeCopyVar(action bot.Action) error {
	from, ok := action.Args["from"].(string)
	if !ok {
		return fmt.Errorf("copy_var requires 'from' argument")
	}
	to, ok := action.Args["to"].(string)
	if !ok {
		return fmt.Errorf("copy_var requires 'to' argument")
	}

	value, exists := ex.mutator.GetVariable(from)
	if !exists {
		return fmt.Errorf("copy_var: source variable '%s' not found", from)
	}

	ex.mutator.SetVariable(to, value)
	return nil
}

func (ex *Executor) executeIncrementVar(action bot.Action) error {
	varName, ok := action.Args["name"].(string)
	if !ok {
		return fmt.Errorf("increment_var requires 'name' argument")
	}

	by := 1
	if val, ok := action.Args["by"]; ok {
		switch v := val.(type) {
		case int:
			by = v
		case float64:
			by = int(v)
		default:
			return fmt.Errorf("increment_var 'by' must be a number")
		}
	}

	return ex.mutator.IncrementVariable(varName, by)
}

func (ex *Executor) executeDeleteVar(action bot.Action) error {
	varName, ok := action.Args["name"].(string)
	if !ok {
		return fmt.Errorf("delete_var requires 'name' argument")
	}

	ex.mutator.DeleteVariable(varName)
	return nil
}
