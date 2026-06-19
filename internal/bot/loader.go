package bot

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFromFile loads a bot definition from a YAML file
func LoadFromFile(path string) (*Bot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read bot file: %w", err)
	}

	var botDef struct {
		Bot   Bot              `yaml:"bot"`
		Flows map[string]*Node `yaml:"flows"`
	}

	if err := yaml.Unmarshal(data, &botDef); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	b := &Bot{
		Name:          botDef.Bot.Name,
		GlobalIntents: botDef.Bot.GlobalIntents,
		Flows:         botDef.Flows,
	}

	if err := b.ValidateBasic(); err != nil {
		return nil, err
	}

	return b, nil
}
