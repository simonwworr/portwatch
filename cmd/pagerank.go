package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var historyDir string
	var format string
	var iterations int
	var damping float64

	cmd := &cobra.Command{
		Use:   "pagerank",
		Short: "Rank hosts by network importance using shared open-port connections",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPageRank(historyDir, format, iterations, damping)
		},
	}

	cmd.Flags().StringVar(&historyDir, "history-dir", ".portwatch/history", "directory containing history files")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	cmd.Flags().IntVar(&iterations, "iterations", 30, "number of PageRank iterations")
	cmd.Flags().Float64Var(&damping, "damping", 0.85, "damping factor (0–1)")

	rootCmd.AddCommand(cmd)
}

func runPageRank(historyDir, format string, iterations int, damping float64) error {
	store, err := history.Load(historyDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return fmt.Errorf("load history: %w", err)
	}

	results := history.PageRank(store, iterations, damping)
	if len(results) == 0 {
		fmt.Println("no data")
		return nil
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(results)
	default:
		fmt.Printf("%-30s  %s\n", "HOST", "SCORE")
		for _, r := range results {
			fmt.Printf("%-30s  %.6f\n", r.Host, r.Score)
		}
	}
	return nil
}
