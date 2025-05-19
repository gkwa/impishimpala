package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version information for impishimpala`,
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion()
		fmt.Printf("impishimpala %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	version := info.Main.Version
	if version == "(devel)" {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value[:8] // short commit hash
			}
		}
		return "development"
	}

	return version
}
