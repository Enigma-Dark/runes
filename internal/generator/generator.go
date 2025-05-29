package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Enigma-Dark/runes/internal/templates"
	"github.com/Enigma-Dark/runes/internal/types"
)

const DefaultTemplate = "enigmadark"

// GenerateConfig holds configuration for test generation
type GenerateConfig struct {
	ContractName string
	OutputFile   string
	ReplayGroups []types.ReplayGroup
	Template     string // Template name to use
}

// templateData holds data for the template
type templateData struct {
	ContractName string
	ReplayGroups []templateReplayGroup
}

// templateReplayGroup represents a replay group with template-formatted calls
type templateReplayGroup struct {
	TestName      string
	TemplateCalls []templateCall
}

// templateCall represents a call in the template
type templateCall struct {
	// Function call fields
	IsFunctionCall bool
	FunctionName   string
	ParamList      string

	// Actor setup fields
	IsSetUpActor bool
	ActorAddress string

	// Delay fields
	IsDelay    bool
	DelayValue string
}

// GenerateFoundryTest generates a Foundry test file from replay groups
func GenerateFoundryTest(config GenerateConfig) error {
	if len(config.ReplayGroups) == 0 {
		return fmt.Errorf("no replay groups to generate")
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(config.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Initialize template manager
	templateManager := templates.NewManager()
	if err := templateManager.LoadBuiltinTemplates(); err != nil {
		return fmt.Errorf("failed to load builtin templates: %w", err)
	}

	// Determine which template to use
	templateName := config.Template
	if templateName == "" {
		templateName = DefaultTemplate
	}

	// Try to load external template if it looks like a file path
	if strings.Contains(templateName, "/") || strings.Contains(templateName, "\\") || strings.HasSuffix(templateName, ".tmpl") {
		baseName := strings.TrimSuffix(filepath.Base(templateName), ".tmpl")
		if err := templateManager.LoadExternalTemplate(baseName, templateName); err != nil {
			return fmt.Errorf("failed to load external template: %w", err)
		}
		templateName = baseName
	}

	// Get the template
	tmpl, err := templateManager.GetTemplate(templateName)
	if err != nil {
		return err
	}

	// Prepare template data
	data := templateData{
		ContractName: config.ContractName,
	}

	// Convert replay groups to template format
	for _, group := range config.ReplayGroups {
		templateGroup := templateReplayGroup{
			TestName:      group.TestName,
			TemplateCalls: convertToTemplateCalls(group.Calls),
		}
		data.ReplayGroups = append(data.ReplayGroups, templateGroup)
	}

	// Create output file
	file, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// ListAvailableTemplates returns a list of available template names
func ListAvailableTemplates() ([]string, error) {
	templateManager := templates.NewManager()
	if err := templateManager.LoadBuiltinTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load builtin templates: %w", err)
	}
	return templateManager.ListTemplates(), nil
}

// convertToTemplateCalls converts ParsedCalls to templateCalls with proper sequencing
func convertToTemplateCalls(calls []types.ParsedCall) []templateCall {
	var result []templateCall
	var lastActor string

	for _, call := range calls {
		// Check if we need to set up a new actor
		currentActor := mapAddressToActor(call.Src)
		if currentActor != lastActor {
			result = append(result, templateCall{
				IsSetUpActor: true,
				ActorAddress: currentActor,
			})
			lastActor = currentActor
		}

		// Add delay if specified
		if call.HasDelay {
			result = append(result, templateCall{
				IsDelay:    true,
				DelayValue: call.DelayValue,
			})
		}

		// Add the function call
		if call.FunctionName != "" {
			var paramValues []string
			for _, param := range call.Parameters {
				paramValues = append(paramValues, param.Value)
			}

			result = append(result, templateCall{
				IsFunctionCall: true,
				FunctionName:   call.FunctionName,
				ParamList:      strings.Join(paramValues, ", "),
			})
		}
	}

	return result
}

// mapAddressToActor maps source addresses to actor constants
func mapAddressToActor(address string) string {
	switch address {
	case "0x0000000000000000000000000000000000010000":
		return "USER1"
	case "0x0000000000000000000000000000000000020000":
		return "USER2"
	case "0x0000000000000000000000000000000000030000":
		return "USER3"
	default:
		// For unknown addresses, map to USER1 as default
		return "USER1"
	}
}
