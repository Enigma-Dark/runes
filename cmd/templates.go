package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Enigma-Dark/runes/internal/generator"
)

// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available templates",
	Long: `List all available builtin templates that can be used for generating Foundry tests.

You can also use custom templates by providing a path to a .tmpl file with the --template flag.

Example:
  runes templates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := generator.ListAvailableTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}

		fmt.Println("Available builtin templates:")
		for _, tmpl := range templates {
			fmt.Printf("  - %s\n", tmpl)
		}

		fmt.Println("\nUsage:")
		fmt.Println("  --template basic        # Use basic template")
		fmt.Println("  --template enigmadark   # Use enigmadark template (default)")
		fmt.Println("  --template /path/to/custom.tmpl  # Use custom template file")

		fmt.Printf("\nCurrent default: %s\n", strings.Join(templates, ", "))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
