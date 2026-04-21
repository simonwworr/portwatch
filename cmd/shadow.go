package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var dir string
	var minScans int
	var maxAppearances int
	var format string

	cmd := &cobra.Command{
		Use:   "shadow",
		Short: "Detect transient (shadow) ports that appear only briefly",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShadow(dir, minScans, maxAppearances, format)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "Directory containing history files")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "Minimum scans required to evaluate a host")
	cmd.Flags().IntVar(&maxAppearances, "max-appearances", 1, "Maximum appearances to classify a port as shadow")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runShadow(dir string, minScans, maxAppearances int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return err
	}

	results := history.Shadow(store, minScans, maxAppearances)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No shadow ports detected.")
		return nil
	}

	fmt.Printf("%-20s %8s %12s %s\n", "HOST", "PORT", "APPEARANCES", "FIRST SEEN")
	for _, r := range results {
		fmt.Printf("%-20s %8d %12d %s\n",
			r.Host, r.Port, r.Appearances, r.FirstSeen.Format("2006-01-02 15:04"))
	}
	return nil
}
