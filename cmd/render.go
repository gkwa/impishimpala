package cmd

import (
	"fmt"

	"github.com/gkwa/impishimpala/internal/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	templatePath string
	renderOutput string
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a boilerplate template",
	Long: `Render a boilerplate template using its configuration and variables,
outputting the result to the specified directory.`,
	RunE: runRender,
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVarP(&templatePath, "template", "t", "", "path to template directory (required)")
	renderCmd.Flags().StringVarP(&renderOutput, "output", "o", "", "output directory for rendered template (required)")
	renderCmd.MarkFlagRequired("template")
	renderCmd.MarkFlagRequired("output")
}

func runRender(cmd *cobra.Command, args []string) error {
	verbosity := viper.GetInt("verbose")

	r := renderer.New(verbosity)

	err := r.RenderTemplate(templatePath, renderOutput)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}
