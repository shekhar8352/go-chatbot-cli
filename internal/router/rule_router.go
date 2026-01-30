package router

import (
	"chatbot-go/internal/bot"
	"strings"
)

// RuleRouter routes based on exact matches, keywords, and simple similarity
type RuleRouter struct{}

// NewRuleRouter creates a new rule-based router
func NewRuleRouter() *RuleRouter {
	return &RuleRouter{}
}

// Route attempts to match user input to an intent using rule-based matching
func (r *RuleRouter) Route(input string, intents []bot.Intent) (string, error) {
	if len(intents) == 0 {
		return "", ErrNoIntents{}
	}

	inputLower := strings.ToLower(strings.TrimSpace(input))

	// 1. Exact match
	for _, intent := range intents {
		for _, example := range intent.Examples {
			if strings.ToLower(example) == inputLower {
				return intent.Name, nil
			}
		}
	}

	// 2. Keyword/substring match
	for _, intent := range intents {
		for _, example := range intent.Examples {
			exampleLower := strings.ToLower(example)
			// Check if input contains example or example contains input
			if strings.Contains(inputLower, exampleLower) || strings.Contains(exampleLower, inputLower) {
				return intent.Name, nil
			}
		}
	}

	// 3. Word-level matching (simple similarity)
	// Split input into words and check if any intent example contains those words
	inputWords := strings.Fields(inputLower)
	for _, intent := range intents {
		for _, example := range intent.Examples {
			exampleLower := strings.ToLower(example)
			exampleWords := strings.Fields(exampleLower)

			// Check if significant words match
			matchCount := 0
			for _, inputWord := range inputWords {
				for _, exampleWord := range exampleWords {
					if inputWord == exampleWord && len(inputWord) > 2 { // Ignore very short words
						matchCount++
						break
					}
				}
			}

			// If at least 50% of words match, consider it a match
			if len(inputWords) > 0 && float64(matchCount)/float64(len(inputWords)) >= 0.5 {
				return intent.Name, nil
			}
		}
	}

	return "", ErrNoMatch{}
}

// ErrNoMatch indicates no intent matched
type ErrNoMatch struct{}

func (e ErrNoMatch) Error() string {
	return "no matching intent found"
}

// ErrNoIntents indicates no intents were provided
type ErrNoIntents struct{}

func (e ErrNoIntents) Error() string {
	return "no intents provided"
}
