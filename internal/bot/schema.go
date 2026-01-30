package bot

// Schema validation helpers
// Note: Full validation is done in internal/validate/flow.go

// ValidateBasic performs basic structural validation
func (b *Bot) ValidateBasic() error {
	if b.Name == "" {
		return ErrInvalidBot("bot name is required")
	}
	if len(b.Flows) == 0 {
		return ErrInvalidBot("at least one flow node is required")
	}
	if _, exists := b.Flows["start"]; !exists {
		return ErrInvalidBot("flow must contain a 'start' node")
	}
	return nil
}

// ErrInvalidBot represents a bot validation error
type ErrInvalidBot string

func (e ErrInvalidBot) Error() string {
	return string(e)
}
