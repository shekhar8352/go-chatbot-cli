package bot

// Bot represents the complete bot definition loaded from YAML
type Bot struct {
	Name          string           `yaml:"name"`
	GlobalIntents []Intent         `yaml:"global_intents,omitempty"`
	Flows         map[string]*Node `yaml:"flows"`
}

// Node represents a single conversation node in the flow
type Node struct {
	Message  string   `yaml:"message"`
	Intents  []Intent `yaml:"intents,omitempty"`
	Input    *Input   `yaml:"input,omitempty"`
	Actions  []Action `yaml:"actions,omitempty"`
	Next     string   `yaml:"next,omitempty"`
	Terminal bool     `yaml:"terminal,omitempty"`
}

// Intent defines an intent that can be matched from user input
type Intent struct {
	Name     string   `yaml:"name"`
	Examples []string `yaml:"examples"`
	Next     string   `yaml:"next"`
}

// Input defines how to capture user input
type Input struct {
	Type         string   `yaml:"type"` // text, choice, confirm, email, number, regex
	SaveAs       string   `yaml:"save_as"`
	Options      []string `yaml:"options,omitempty"`
	Pattern      string   `yaml:"pattern,omitempty"`
	RetryMessage string   `yaml:"retry_message,omitempty"`
	OnYes        string   `yaml:"on_yes,omitempty"`
	OnNo         string   `yaml:"on_no,omitempty"`
}

// Action represents an action to execute
type Action struct {
	Type string                 `yaml:"type"`
	Args map[string]interface{} `yaml:"args,omitempty"`
}
