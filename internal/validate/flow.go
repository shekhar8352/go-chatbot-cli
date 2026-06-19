package validate

import (
	"chatbot-go/internal/bot"
	"fmt"
)

// ValidateFlow performs comprehensive flow validation
func ValidateFlow(b *bot.Bot) error {
	if err := b.ValidateBasic(); err != nil {
		return err
	}

	for nodeName, node := range b.Flows {
		if node.Input != nil && len(node.Intents) > 0 {
			return fmt.Errorf("node '%s' has both input and intents; input takes priority", nodeName)
		}

		if err := bot.ValidateInput(nodeName, node.Input); err != nil {
			return err
		}

		for _, action := range node.Actions {
			if err := bot.ValidateAction(nodeName, action); err != nil {
				return err
			}
		}

		if node.Next != "" {
			if _, exists := b.Flows[node.Next]; !exists {
				return fmt.Errorf("node '%s' references non-existent next node '%s'", nodeName, node.Next)
			}
		}

		if node.Input != nil && node.Input.Type == "confirm" {
			for _, target := range []string{node.Input.OnYes, node.Input.OnNo} {
				if _, exists := b.Flows[target]; !exists {
					return fmt.Errorf("node '%s' confirm input references non-existent node '%s'", nodeName, target)
				}
			}
		}

		for _, intent := range node.Intents {
			if intent.Next != "" {
				if _, exists := b.Flows[intent.Next]; !exists {
					return fmt.Errorf("node '%s' intent '%s' references non-existent next node '%s'", nodeName, intent.Name, intent.Next)
				}
			}
		}
	}

	for _, intent := range b.GlobalIntents {
		if intent.Next != "" {
			if _, exists := b.Flows[intent.Next]; !exists {
				return fmt.Errorf("global intent '%s' references non-existent next node '%s'", intent.Name, intent.Next)
			}
		}
	}

	return nil
}
