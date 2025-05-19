package models

import (
	"github.com/gruntwork-io/boilerplate/config"
)

// TemplateManifest represents a collection of discovered boilerplate templates
type TemplateManifest struct {
	Templates []Template `json:"templates"`
	Version   string     `json:"version"`
}

// Template represents a single boilerplate template with its configuration and files
type Template struct {
	Name             string                    `json:"name"`
	Path             string                    `json:"path"`
	ConfigPath       string                    `json:"config_path"`
	Config           *config.BoilerplateConfig `json:"config,omitempty"`
	TemplateFiles    []TemplateFile            `json:"template_files"`
	HasValidConfig   bool                      `json:"has_valid_config"`
	ConfigParseError string                    `json:"config_parse_error,omitempty"`
}

// TemplateFile represents a single template file within a template directory
type TemplateFile struct {
	RelativePath string `json:"relative_path"`
	AbsolutePath string `json:"absolute_path"`
	Content      string `json:"content,omitempty"`
	IsText       bool   `json:"is_text"`
	Size         int64  `json:"size"`
}
