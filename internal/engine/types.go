package engine

import (
	"chatbot-go/internal/bot"
	"fmt"
)

// Session represents the current conversation session
type Session struct {
	CurrentNode string
	Variables   map[string]string
	History     []Turn
}

// GetVariables returns all variables (implements render.SessionView)
func (s *Session) GetVariables() map[string]string {
	if s.Variables == nil {
		return make(map[string]string)
	}
	return s.Variables
}

// Turn represents a single turn in the conversation
type Turn struct {
	Node      string
	UserInput string
	Response  string
}

// Engine manages the conversation flow using FSM
type Engine struct {
	bot     *bot.Bot
	session *Session
}

// NewEngine creates a new conversation engine
func NewEngine(b *bot.Bot) *Engine {
	return &Engine{
		bot: b,
		session: &Session{
			CurrentNode: "start",
			Variables:   make(map[string]string),
			History:     []Turn{},
		},
	}
}

// GetSession returns the current session
func (e *Engine) GetSession() *Session {
	return e.session
}

// GetCurrentNode returns the current node
func (e *Engine) GetCurrentNode() (*bot.Node, error) {
	node, exists := e.bot.Flows[e.session.CurrentNode]
	if !exists {
		return nil, ErrNodeNotFound(e.session.CurrentNode)
	}
	return node, nil
}

// ErrNodeNotFound represents a node not found error
type ErrNodeNotFound string

func (e ErrNodeNotFound) Error() string {
	return fmt.Sprintf("node '%s' not found", string(e))
}
