package validate

import (
	"chatbot-go/internal/bot"
	"fmt"
)

// ValidateFlow performs comprehensive flow validation
func ValidateFlow(b *bot.Bot) error {
	// Basic validation
	if err := b.ValidateBasic(); err != nil {
		return err
	}

	// Check all referenced nodes exist
	for nodeName, node := range b.Flows {
		// Check next node
		if node.Next != "" {
			if _, exists := b.Flows[node.Next]; !exists {
				return fmt.Errorf("node '%s' references non-existent next node '%s'", nodeName, node.Next)
			}
		}

		// Check intent next nodes
		for _, intent := range node.Intents {
			if intent.Next != "" {
				if _, exists := b.Flows[intent.Next]; !exists {
					return fmt.Errorf("node '%s' intent '%s' references non-existent next node '%s'", nodeName, intent.Name, intent.Next)
				}
			}
		}
	}

	return nil
}
