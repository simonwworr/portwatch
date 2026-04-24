package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var (
		dir       string
		minScans  int
		formatOut string
	)

	cmd := &cobra.Command{
		Use:   "stability",
		Short: "Show port-set stability score per host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStability(dir, minScans, formatOut)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "history", "History directory")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "Minimum scans required")
	cmd.Flags().StringVar(&formatOut, "format", "text", "Output format: text|json")

	rootCmd.AddCommand(cmd)
}

func runStability(dir string, minScans int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return fmt.Errorf("load history: %w", err)
	}

	results := history.Stability(store, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No stability data available.")
		return nil
	}

	fmt.Printf("%-20s %6s %8s %10s  %s\n", "HOST", "SCANS", "CHANGES", "STABILITY", "STABLE PORTS")
	for _, r := range results {
		fmt.Printf("%-20s %6d %8d %9.1f%%  %v\n",
			r.Host, r.Scans, r.Changes, r.Stability*100, r.StablePorts)
	}
	return nil
}
