package renderer

import (
	"fmt"
	"log"

	"github.com/gruntwork-io/boilerplate/options"
	"github.com/gruntwork-io/boilerplate/templates"
	"github.com/gruntwork-io/boilerplate/variables"
)

// Renderer handles rendering of boilerplate templates
type Renderer struct {
	verbosity int
}

// New creates a new Renderer instance
func New(verbosity int) *Renderer {
	return &Renderer{
		verbosity: verbosity,
	}
}

// RenderTemplate renders a boilerplate template to the specified output directory
func (r *Renderer) RenderTemplate(templatePath, outputPath string) error {
	if r.verbosity > 0 {
		log.Printf("Rendering template from %s to %s", templatePath, outputPath)
	}

	// Create boilerplate options
	opts := &options.BoilerplateOptions{
		TemplateUrl:     templatePath,
		TemplateFolder:  templatePath,
		OutputFolder:    outputPath,
		NonInteractive:  true,
		OnMissingKey:    options.ExitWithError,
		OnMissingConfig: options.Exit,
		Vars:            make(map[string]interface{}),
	}

	// Create an empty dependency for the root template
	emptyDep := variables.Dependency{}

	// Process the template
	err := templates.ProcessTemplate(opts, opts, emptyDep)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	if r.verbosity > 0 {
		log.Printf("Successfully rendered template to %s", outputPath)
	}

	return nil
}
