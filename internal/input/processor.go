package input

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"chatbot-go/internal/bot"
)

// Result holds the outcome of processing user input.
type Result struct {
	Value    string
	NextNode string
}

// Processor validates and normalizes captured input.
type Processor struct{}

// NewProcessor creates a new input processor.
func NewProcessor() *Processor {
	return &Processor{}
}

// Process validates raw input against the input definition.
func (p *Processor) Process(def *bot.Input, raw string) (Result, error) {
	inputType := def.Type
	if inputType == "" {
		inputType = "text"
	}

	raw = strings.TrimSpace(raw)

	switch inputType {
	case "text":
		return Result{Value: raw}, nil

	case "choice":
		value, ok := matchChoice(raw, def.Options)
		if !ok {
			return Result{}, fmt.Errorf("invalid choice")
		}
		return Result{Value: value}, nil

	case "confirm":
		value, ok := matchConfirm(raw)
		if !ok {
			return Result{}, fmt.Errorf("invalid confirm")
		}
		next := def.OnNo
		if value == "yes" {
			next = def.OnYes
		}
		return Result{Value: value, NextNode: next}, nil

	case "email":
		if !emailPattern.MatchString(raw) {
			return Result{}, fmt.Errorf("invalid email")
		}
		return Result{Value: strings.ToLower(raw)}, nil

	case "number":
		if _, err := strconv.ParseFloat(raw, 64); err != nil {
			return Result{}, fmt.Errorf("invalid number")
		}
		return Result{Value: raw}, nil

	case "regex":
		re, err := regexp.Compile(def.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("invalid regex pattern: %w", err)
		}
		if !re.MatchString(raw) {
			return Result{}, fmt.Errorf("input does not match pattern")
		}
		return Result{Value: raw}, nil

	default:
		return Result{}, fmt.Errorf("unknown input type: %s", inputType)
	}
}

// RetryMessage returns the message shown when validation fails.
func RetryMessage(def *bot.Input) string {
	if def.RetryMessage != "" {
		return def.RetryMessage
	}

	switch def.Type {
	case "choice":
		return fmt.Sprintf("Please choose one of: %s", strings.Join(def.Options, ", "))
	case "confirm":
		return "Please answer yes or no."
	case "email":
		return "Please enter a valid email address."
	case "number":
		return "Please enter a valid number."
	case "regex":
		return "That input is not valid. Please try again."
	default:
		return "Invalid input. Please try again."
	}
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func matchChoice(raw string, options []string) (string, bool) {
	rawLower := strings.ToLower(raw)

	for _, opt := range options {
		if strings.ToLower(opt) == rawLower {
			return opt, true
		}
	}

	if idx, err := strconv.Atoi(raw); err == nil && idx >= 1 && idx <= len(options) {
		return options[idx-1], true
	}

	for _, opt := range options {
		optLower := strings.ToLower(opt)
		if strings.Contains(rawLower, optLower) || strings.Contains(optLower, rawLower) {
			return opt, true
		}
	}

	return "", false
}

func matchConfirm(raw string) (string, bool) {
	rawLower := strings.ToLower(strings.TrimSpace(raw))

	yesWords := []string{"yes", "y", "yeah", "yep", "confirm", "ok", "sure", "true", "1"}
	noWords := []string{"no", "n", "nope", "cancel", "false", "0"}

	for _, w := range yesWords {
		if rawLower == w {
			return "yes", true
		}
	}
	for _, w := range noWords {
		if rawLower == w {
			return "no", true
		}
	}

	return "", false
}
