package router

import "chatbot-go/internal/bot"

// Router interface for routing user input to intents
type Router interface {
	Route(input string, intents []bot.Intent) (string, error)
}

// RouteResult represents the result of routing
type RouteResult struct {
	IntentName string
	Confidence float64
}
