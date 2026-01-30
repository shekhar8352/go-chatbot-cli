package bot

// Bot represents the complete bot definition loaded from YAML
type Bot struct {
	Name  string           `yaml:"name"`
	Flows map[string]*Node `yaml:"flows"`
}

// Node represents a single conversation node in the flow
type Node struct {
	Message string   `yaml:"message"`
	Intents []Intent `yaml:"intents,omitempty"`
	Input   *Input   `yaml:"input,omitempty"`
	Actions []Action `yaml:"actions,omitempty"`
	Next    string   `yaml:"next,omitempty"`
}

// Intent defines an intent that can be matched from user input
type Intent struct {
	Name     string   `yaml:"name"`
	Examples []string `yaml:"examples"`
	Next     string   `yaml:"next"`
}

// Input defines how to capture user input
type Input struct {
	Type   string `yaml:"type"` // "text" for now
	SaveAs string `yaml:"save_as"`
}

// Action represents an action to execute
type Action struct {
	Type string                 `yaml:"type"` // "set_var" for now
	Args map[string]interface{} `yaml:"args,omitempty"`
}
