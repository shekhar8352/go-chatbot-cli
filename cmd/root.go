package cmd

import (
	"context"
	"fmt"
	"os"

	"chatbot-go/internal/bot"
	"chatbot-go/internal/engine"
	"chatbot-go/internal/llm"
	"chatbot-go/internal/validate"

	"github.com/spf13/cobra"
)

var (
	botFile     string
	llmType     string
	ollamaURL   string
	ollamaModel string
)

var rootCmd = &cobra.Command{
	Use:   "chatbot",
	Short: "A deterministic CLI chatbot framework",
	Long: `A CLI-based chatbot framework that supports deterministic conversational flows
and is LLM-ready for integration with local LLMs.`,
	RunE: runChatbot,
}

func init() {
	rootCmd.Flags().StringVarP(&botFile, "bot", "b", "examples/support-bot.yaml", "Path to bot YAML file")
	rootCmd.Flags().StringVarP(&llmType, "llm", "l", "noop", "LLM provider type (noop, ollama)")
	rootCmd.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "Ollama API URL")
	rootCmd.Flags().StringVar(&ollamaModel, "ollama-model", "llama2", "Ollama model name")
}

func runChatbot(cmd *cobra.Command, args []string) error {
	// Load bot
	b, err := bot.LoadFromFile(botFile)
	if err != nil {
		return fmt.Errorf("failed to load bot: %w", err)
	}

	// Validate flow
	if err := validate.ValidateFlow(b); err != nil {
		return fmt.Errorf("flow validation failed: %w", err)
	}

	// Initialize LLM provider
	var llmProvider llm.Provider
	switch llmType {
	case "ollama":
		llmProvider = llm.NewOllamaProvider(ollamaURL, ollamaModel)
	case "noop", "":
		llmProvider = llm.NewNoopProvider()
	default:
		return fmt.Errorf("unknown LLM provider: %s", llmType)
	}

	// Create and run engine
	conversationEngine := engine.NewConversationEngine(b, llmProvider)
	ctx := context.Background()

	if err := conversationEngine.Run(ctx); err != nil {
		return fmt.Errorf("conversation error: %w", err)
	}

	return nil
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
