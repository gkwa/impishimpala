package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gkwa/impishimpala/internal/gatherer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputFile string

var gatherCmd = &cobra.Command{
	Use:   "gather [paths...]",
	Short: "Gather boilerplate templates from directories",
	Long: `Recursively discover boilerplate templates and their configuration files
from one or more directories, creating a manifest of all discovered templates.`,
	Args: cobra.MinimumNArgs(1),
	RunE: runGather,
}

func init() {
	rootCmd.AddCommand(gatherCmd)
	gatherCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output manifest to file (default: stdout)")
}

func runGather(cmd *cobra.Command, args []string) error {
	verbosity := viper.GetInt("verbose")

	g := gatherer.New(verbosity)

	manifest, err := g.GatherTemplates(args)
	if err != nil {
		return fmt.Errorf("failed to gather templates: %w", err)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Output to file or stdout
	if outputFile != "" {
		if verbosity > 0 {
			log.Printf("Writing manifest to %s", outputFile)
		}
		return os.WriteFile(outputFile, jsonData, 0o644)
	}

	fmt.Print(string(jsonData))
	return nil
}
