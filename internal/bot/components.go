package bot

import "fmt"

// Supported input types.
var SupportedInputTypes = map[string]bool{
	"text":    true,
	"choice":  true,
	"confirm": true,
	"email":   true,
	"number":  true,
	"regex":   true,
}

// Supported action types.
var SupportedActionTypes = map[string]bool{
	"set_var":      true,
	"append_var":   true,
	"copy_var":     true,
	"increment_var": true,
	"delete_var":   true,
}

// ValidateInput validates an input component definition.
func ValidateInput(nodeName string, input *Input) error {
	if input == nil {
		return nil
	}

	if input.SaveAs == "" {
		return fmt.Errorf("node '%s': input requires save_as", nodeName)
	}

	inputType := input.Type
	if inputType == "" {
		inputType = "text"
	}
	if !SupportedInputTypes[inputType] {
		return fmt.Errorf("node '%s': unknown input type '%s'", nodeName, inputType)
	}

	switch inputType {
	case "choice":
		if len(input.Options) == 0 {
			return fmt.Errorf("node '%s': choice input requires options", nodeName)
		}
	case "confirm":
		if input.OnYes == "" || input.OnNo == "" {
			return fmt.Errorf("node '%s': confirm input requires on_yes and on_no", nodeName)
		}
	case "regex":
		if input.Pattern == "" {
			return fmt.Errorf("node '%s': regex input requires pattern", nodeName)
		}
	}

	return nil
}

// ValidateAction validates an action component definition.
func ValidateAction(nodeName string, action Action) error {
	if !SupportedActionTypes[action.Type] {
		return fmt.Errorf("node '%s': unknown action type '%s'", nodeName, action.Type)
	}

	switch action.Type {
	case "set_var", "append_var", "delete_var", "increment_var":
		if _, ok := action.Args["name"].(string); !ok {
			return fmt.Errorf("node '%s': action '%s' requires args.name", nodeName, action.Type)
		}
	case "copy_var":
		if _, ok := action.Args["from"].(string); !ok {
			return fmt.Errorf("node '%s': copy_var requires args.from", nodeName)
		}
		if _, ok := action.Args["to"].(string); !ok {
			return fmt.Errorf("node '%s': copy_var requires args.to", nodeName)
		}
	}

	return nil
}
