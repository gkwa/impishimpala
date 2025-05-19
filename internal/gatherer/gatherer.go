package gatherer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gkwa/impishimpala/internal/models"
	"github.com/gkwa/impishimpala/internal/utils"
	"github.com/gruntwork-io/boilerplate/config"
	"github.com/gruntwork-io/boilerplate/options"
	"github.com/gruntwork-io/boilerplate/util"
)

// Gatherer handles discovery and collection of boilerplate templates
type Gatherer struct {
	verbosity int
}

// New creates a new Gatherer instance
func New(verbosity int) *Gatherer {
	return &Gatherer{
		verbosity: verbosity,
	}
}

// GatherTemplates recursively discovers boilerplate templates in the given paths
func (g *Gatherer) GatherTemplates(paths []string) (*models.TemplateManifest, error) {
	// Redirect boilerplate library's internal logging to stderr
	util.Logger.SetOutput(os.Stderr)

	manifest := &models.TemplateManifest{
		Templates: make([]models.Template, 0),
		Version:   "1.0",
	}

	for _, rootPath := range paths {
		if g.verbosity > 0 {
			log.Printf("Scanning directory: %s", rootPath)
		}

		err := g.scanDirectory(rootPath, manifest)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory %s: %w", rootPath, err)
		}
	}

	if g.verbosity > 0 {
		log.Printf("Found %d templates", len(manifest.Templates))
	}

	return manifest, nil
}

// scanDirectory recursively scans a directory for boilerplate templates
func (g *Gatherer) scanDirectory(rootPath string, manifest *models.TemplateManifest) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if g.verbosity > 1 {
				log.Printf("Warning: failed to access %s: %v", path, err)
			}
			return nil // Continue walking despite errors
		}

		// Check if this is a boilerplate.yml file
		if info.Name() == "boilerplate.yml" && !info.IsDir() {
			templateDir := filepath.Dir(path)
			template, err := g.processTemplate(templateDir, path)
			if err != nil {
				if g.verbosity > 0 {
					log.Printf("Warning: failed to process template at %s: %v", templateDir, err)
				}
				return nil // Continue despite processing errors
			}

			manifest.Templates = append(manifest.Templates, *template)
			if g.verbosity > 1 {
				log.Printf("Discovered template: %s", template.Name)
			}
		}

		return nil
	})
}

// processTemplate processes a single template directory
func (g *Gatherer) processTemplate(templateDir, configPath string) (*models.Template, error) {
	absTemplatePath, err := filepath.Abs(templateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute config path: %w", err)
	}

	template := &models.Template{
		Name:          filepath.Base(absTemplatePath),
		Path:          absTemplatePath,
		ConfigPath:    absConfigPath,
		TemplateFiles: make([]models.TemplateFile, 0),
	}

	// Parse boilerplate config
	boilerplateConfig, err := config.LoadBoilerplateConfig(&options.BoilerplateOptions{
		TemplateFolder: absTemplatePath,
	})
	if err != nil {
		template.HasValidConfig = false
		template.ConfigParseError = err.Error()
		if g.verbosity > 1 {
			log.Printf("Warning: failed to parse config for %s: %v", template.Name, err)
		}
	} else {
		template.HasValidConfig = true
		template.Config = boilerplateConfig
	}

	// Gather template files
	err = g.gatherTemplateFiles(absTemplatePath, template)
	if err != nil {
		return nil, fmt.Errorf("failed to gather template files: %w", err)
	}

	return template, nil
}

// gatherTemplateFiles collects all files in the template directory
func (g *Gatherer) gatherTemplateFiles(templateDir string, template *models.Template) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		// Skip the boilerplate.yml file itself
		if info.Name() == "boilerplate.yml" {
			return nil
		}

		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return nil // Skip if we can't get relative path
		}

		templateFile := models.TemplateFile{
			RelativePath: relPath,
			AbsolutePath: path,
			Size:         info.Size(),
		}

		// Check if file is text and read content if requested
		isText, err := utils.IsTextFile(path)
		if err != nil {
			if g.verbosity > 2 {
				log.Printf("Warning: failed to determine file type for %s: %v", path, err)
			}
			templateFile.IsText = false
		} else {
			templateFile.IsText = isText

			// Read content for text files (with size limit)
			if isText && info.Size() < 1024*1024 { // 1MB limit
				content, err := os.ReadFile(path)
				if err != nil {
					if g.verbosity > 2 {
						log.Printf("Warning: failed to read file %s: %v", path, err)
					}
				} else {
					templateFile.Content = string(content)
				}
			}
		}

		template.TemplateFiles = append(template.TemplateFiles, templateFile)
		return nil
	})
}
