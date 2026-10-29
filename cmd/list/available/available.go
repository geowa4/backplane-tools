package available

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/openshift/backplane-tools/pkg/tools"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var output string
	availableCmd := &cobra.Command{
		Use:     "available",
		Args:    cobra.NoArgs,
		Aliases: []string{"installable", "possible"},
		Short:   "List available tools for install",
		Long:    "List tools that are available to install with backplane-tools",
		RunE: func(_ *cobra.Command, _ []string) error {
			switch output {
			case "text":
				return List()
			case "json":
				return ListJSON()
			default:
				return fmt.Errorf("unsupported output format '%s': must be one of 'text' or 'json'", output)
			}
		},
	}
	availableCmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: 'text' or 'json'")
	return availableCmd
}

func List() error {
	fmt.Println("The following tools are available for install:")

	toolMap := tools.GetMap()
	for _, t := range toolMap {
		version, err := t.LatestVersion()
		if err != nil {
			return fmt.Errorf("failed to determine version for '%s': %w", t.Name(), err)
		}
		fmt.Printf("- %s %s\n", t.Name(), version)
	}
	return nil
}

// ListJSON prints a JSON object mapping each tool's name to its latest available version
func ListJSON() error {
	versions, err := LatestVersions(tools.GetMap())
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(versions)
}

// LatestVersions resolves the latest available version of each provided tool
func LatestVersions(toolMap map[string]tools.Tool) (map[string]string, error) {
	versions := make(map[string]string, len(toolMap))
	for name, t := range toolMap {
		version, err := t.LatestVersion()
		if err != nil {
			return nil, fmt.Errorf("failed to determine version for '%s': %w", name, err)
		}
		versions[name] = version
	}
	return versions, nil
}
