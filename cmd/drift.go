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
	var format string

	driftCmd := &cobra.Command{
		Use:   "drift",
		Short: "Show how each host's open ports have drifted from their first recorded snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDrift(dir, format)
		},
	}

	driftCmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	driftCmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(driftCmd)
}

func runDrift(dir, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Drift(store)
	if len(results) == 0 {
		fmt.Println("no drift data available (need at least 2 snapshots per host)")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	for _, r := range results {
		fmt.Printf("host: %s  drift_score: %.2f\n", r.Host, r.DriftScore)
		if len(r.Added) > 0 {
			fmt.Printf("  added:   %v\n", r.Added)
		}
		if len(r.Removed) > 0 {
			fmt.Printf("  removed: %v\n", r.Removed)
		}
		if len(r.Stable) > 0 {
			fmt.Printf("  stable:  %v\n", r.Stable)
		}
	}
	return nil
}
