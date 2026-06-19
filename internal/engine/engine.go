package engine

import (
	"context"
	"fmt"

	"chatbot-go/internal/actions"
	"chatbot-go/internal/bot"
	"chatbot-go/internal/input"
	"chatbot-go/internal/llm"
	"chatbot-go/internal/render"
	"chatbot-go/internal/router"
)

// ConversationEngine orchestrates the conversation flow
type ConversationEngine struct {
	engine         *Engine
	ruleRouter     router.Router
	llmRouter      *router.LLMRouter
	llmProvider    llm.Provider
	renderer       *render.CLIRenderer
	executor       *actions.Executor
	inputProcessor *input.Processor
}

// NewConversationEngine creates a new conversation engine
func NewConversationEngine(b *bot.Bot, llmProvider llm.Provider) *ConversationEngine {
	eng := NewEngine(b)
	return &ConversationEngine{
		engine:         eng,
		ruleRouter:     router.NewRuleRouter(),
		llmRouter:      router.NewLLMRouter(llmProvider),
		llmProvider:    llmProvider,
		renderer:       render.NewCLIRenderer(),
		executor:       actions.NewExecutor(eng),
		inputProcessor: input.NewProcessor(),
	}
}

// Run starts the conversation loop
func (ce *ConversationEngine) Run(ctx context.Context) error {
	for {
		node, err := ce.engine.GetCurrentNode()
		if err != nil {
			return fmt.Errorf("failed to get current node: %w", err)
		}

		message := ce.renderer.RenderMessage(node, ce.engine.GetSession())
		ce.renderer.PrintMessage(message)

		isTerminal, err := ce.engine.IsTerminal()
		if err != nil {
			return err
		}
		if isTerminal {
			return nil
		}

		effectiveIntents := ce.engine.EffectiveIntents(node)
		if len(effectiveIntents) > 0 {
			ce.renderer.ShowIntents(effectiveIntents)
		}
		if node.Input != nil && node.Input.Type == "choice" && len(node.Input.Options) > 0 {
			ce.renderer.ShowInputOptions(node.Input.Options)
		}

		userInput, err := ce.renderer.ReadInput()
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		if node.Input != nil {
			if err := ce.handleInput(ctx, node, userInput); err != nil {
				return err
			}
			continue
		}

		if len(effectiveIntents) > 0 {
			if err := ce.handleIntents(ctx, node, effectiveIntents, userInput); err != nil {
				return err
			}
		} else if node.Next != "" {
			for _, action := range node.Actions {
				if err := ce.executor.Execute(action, userInput); err != nil {
					return fmt.Errorf("action execution failed: %w", err)
				}
			}
			if err := ce.engine.Transition(node.Next); err != nil {
				return fmt.Errorf("transition failed: %w", err)
			}
		}

		ce.engine.AddTurn(node.Message, userInput, message)
	}
}

func (ce *ConversationEngine) handleInput(ctx context.Context, node *bot.Node, userInput string) error {
	result, err := ce.inputProcessor.Process(node.Input, userInput)
	if err != nil {
		ce.renderer.PrintMessage(input.RetryMessage(node.Input))
		return nil
	}

	ce.engine.SetVariable(node.Input.SaveAs, result.Value)

	for _, action := range node.Actions {
		if err := ce.executor.Execute(action, result.Value); err != nil {
			return fmt.Errorf("action execution failed: %w", err)
		}
	}

	nextNode := node.Next
	if result.NextNode != "" {
		nextNode = result.NextNode
	}

	if nextNode != "" {
		if err := ce.engine.Transition(nextNode); err != nil {
			return fmt.Errorf("transition failed: %w", err)
		}
	}

	return nil
}

func (ce *ConversationEngine) handleIntents(ctx context.Context, node *bot.Node, intents []bot.Intent, userInput string) error {
	intentName, err := ce.ruleRouter.Route(userInput, intents)
	if err != nil {
		if ce.llmProvider != nil {
			intentName, err = ce.llmRouter.Route(ctx, userInput, intents)
			if err != nil {
				ce.renderer.PrintMessage("I didn't understand that. Please try again.")
				return nil
			}
		} else {
			ce.renderer.PrintMessage("I didn't understand that. Please try again.")
			return nil
		}
	}

	for _, intent := range intents {
		if intent.Name == intentName {
			for _, action := range node.Actions {
				if err := ce.executor.Execute(action, userInput); err != nil {
					return fmt.Errorf("action execution failed: %w", err)
				}
			}
			if intent.Next != "" {
				if err := ce.engine.Transition(intent.Next); err != nil {
					return fmt.Errorf("transition failed: %w", err)
				}
			}
			break
		}
	}

	return nil
}
