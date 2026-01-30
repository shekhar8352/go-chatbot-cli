package render

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chatbot-go/internal/bot"
)

// SessionView defines the interface for accessing session data for rendering
type SessionView interface {
	GetVariables() map[string]string
}

// CLIRenderer renders conversation to CLI
type CLIRenderer struct {
	reader *bufio.Reader
}

// NewCLIRenderer creates a new CLI renderer
func NewCLIRenderer() *CLIRenderer {
	return &CLIRenderer{
		reader: bufio.NewReader(os.Stdin),
	}
}

// RenderMessage renders a node message with variable interpolation
func (r *CLIRenderer) RenderMessage(node *bot.Node, session SessionView) string {
	message := node.Message
	// Interpolate variables: {{var_name}}
	for key, value := range session.GetVariables() {
		placeholder := fmt.Sprintf("{{%s}}", key)
		message = strings.ReplaceAll(message, placeholder, value)
	}
	return message
}

// PrintMessage prints a message to stdout
func (r *CLIRenderer) PrintMessage(message string) {
	fmt.Println(message)
}

// ReadInput reads user input from stdin
func (r *CLIRenderer) ReadInput() (string, error) {
	fmt.Print("> ")
	input, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ShowIntents displays available intent options (if applicable)
func (r *CLIRenderer) ShowIntents(intents []bot.Intent) {
	if len(intents) == 0 {
		return
	}
	fmt.Println("\nAvailable options:")
	for i, intent := range intents {
		if len(intent.Examples) > 0 {
			fmt.Printf("  %d. %s (e.g., \"%s\")\n", i+1, intent.Name, intent.Examples[0])
		} else {
			fmt.Printf("  %d. %s\n", i+1, intent.Name)
		}
	}
	fmt.Println()
}
