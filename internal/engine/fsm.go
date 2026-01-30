package engine

import (
	"fmt"
)

// Transition moves the session to a new node
func (e *Engine) Transition(nextNode string) error {
	if nextNode == "" {
		return ErrInvalidTransition("next node cannot be empty")
	}

	if _, exists := e.bot.Flows[nextNode]; !exists {
		return ErrInvalidTransition(fmt.Sprintf("node '%s' does not exist", nextNode))
	}

	e.session.CurrentNode = nextNode
	return nil
}

// IsTerminal checks if the current node is terminal (no next node and no intents)
func (e *Engine) IsTerminal() (bool, error) {
	node, err := e.GetCurrentNode()
	if err != nil {
		return false, err
	}

	// Terminal if no next node and no intents (or empty intents)
	return node.Next == "" && len(node.Intents) == 0, nil
}

// ErrInvalidTransition represents an invalid state transition error
type ErrInvalidTransition string

func (e ErrInvalidTransition) Error() string {
	return string(e)
}
