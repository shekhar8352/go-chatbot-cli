package engine

// SetVariable sets a variable in the session
func (e *Engine) SetVariable(key, value string) {
	if e.session.Variables == nil {
		e.session.Variables = make(map[string]string)
	}
	e.session.Variables[key] = value
}

// GetVariable retrieves a variable from the session
func (e *Engine) GetVariable(key string) (string, bool) {
	if e.session.Variables == nil {
		return "", false
	}
	val, exists := e.session.Variables[key]
	return val, exists
}

// AddTurn adds a turn to the conversation history
func (e *Engine) AddTurn(node, userInput, response string) {
	e.session.History = append(e.session.History, Turn{
		Node:      node,
		UserInput: userInput,
		Response:  response,
	})
}
