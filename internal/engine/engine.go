package engine

import (
	"context"
	"fmt"

	"chatbot-go/internal/actions"
	"chatbot-go/internal/bot"
	"chatbot-go/internal/llm"
	"chatbot-go/internal/render"
	"chatbot-go/internal/router"
)

// ConversationEngine orchestrates the conversation flow
type ConversationEngine struct {
	engine      *Engine
	ruleRouter  router.Router
	llmRouter   *router.LLMRouter
	llmProvider llm.Provider
	renderer    *render.CLIRenderer
	executor    *actions.Executor
}

// NewConversationEngine creates a new conversation engine
func NewConversationEngine(b *bot.Bot, llmProvider llm.Provider) *ConversationEngine {
	eng := NewEngine(b)
	return &ConversationEngine{
		engine:      eng,
		ruleRouter:  router.NewRuleRouter(),
		llmRouter:   router.NewLLMRouter(llmProvider),
		llmProvider: llmProvider,
		renderer:    render.NewCLIRenderer(),
		executor:    actions.NewExecutor(eng),
	}
}

// Run starts the conversation loop
func (ce *ConversationEngine) Run(ctx context.Context) error {
	for {
		// Get current node
		node, err := ce.engine.GetCurrentNode()
		if err != nil {
			return fmt.Errorf("failed to get current node: %w", err)
		}

		// Render and display message
		message := ce.renderer.RenderMessage(node, ce.engine.GetSession())
		ce.renderer.PrintMessage(message)

		// Check if terminal
		isTerminal, err := ce.engine.IsTerminal()
		if err != nil {
			return err
		}
		if isTerminal {
			return nil
		}

		// Show available intents if any
		if len(node.Intents) > 0 {
			ce.renderer.ShowIntents(node.Intents)
		}

		// Read user input
		userInput, err := ce.renderer.ReadInput()
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		// Handle input capture (if node has input definition)
		if node.Input != nil {
			// Save input directly to variable
			ce.engine.SetVariable(node.Input.SaveAs, userInput)

			// Execute any actions
			for _, action := range node.Actions {
				if err := ce.executor.Execute(action, userInput); err != nil {
					return fmt.Errorf("action execution failed: %w", err)
				}
			}

			// Transition to next node
			if node.Next != "" {
				if err := ce.engine.Transition(node.Next); err != nil {
					return fmt.Errorf("transition failed: %w", err)
				}
			}
			continue
		}

		// Handle intent-based routing
		if len(node.Intents) > 0 {
			// Try rule router first
			intentName, err := ce.ruleRouter.Route(userInput, node.Intents)
			if err != nil {
				// If rule router fails, try LLM router (if available)
				if ce.llmProvider != nil {
					intentName, err = ce.llmRouter.Route(ctx, userInput, node.Intents)
					if err != nil {
						ce.renderer.PrintMessage("I didn't understand that. Please try again.")
						continue
					}
				} else {
					ce.renderer.PrintMessage("I didn't understand that. Please try again.")
					continue
				}
			}

			// Find the matched intent and transition
			for _, intent := range node.Intents {
				if intent.Name == intentName {
					// Execute any actions
					for _, action := range node.Actions {
						if err := ce.executor.Execute(action, userInput); err != nil {
							return fmt.Errorf("action execution failed: %w", err)
						}
					}

					// Transition to next node
					if intent.Next != "" {
						if err := ce.engine.Transition(intent.Next); err != nil {
							return fmt.Errorf("transition failed: %w", err)
						}
					}
					break
				}
			}
		} else if node.Next != "" {
			// No intents, just transition to next
			// Execute any actions first
			for _, action := range node.Actions {
				if err := ce.executor.Execute(action, userInput); err != nil {
					return fmt.Errorf("action execution failed: %w", err)
				}
			}

			if err := ce.engine.Transition(node.Next); err != nil {
				return fmt.Errorf("transition failed: %w", err)
			}
		}

		// Record turn in history
		ce.engine.AddTurn(node.Message, userInput, message)
	}
}
