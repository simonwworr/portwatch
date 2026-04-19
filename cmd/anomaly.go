package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var histDir string
	var threshold float64
	var format string

	cmd := &cobra.Command{
		Use:   "anomaly",
		Short: "Detect ports that appear rarely across scan history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnomaly(histDir, threshold, format)
		},
	}

	cmd.Flags().StringVar(&histDir, "history-dir", "history", "Directory containing history files")
	cmd.Flags().Float64Var(&threshold, "threshold", 50.0, "Percentage below which a port is considered rare")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runAnomaly(histDir string, threshold float64, format string) error {
	s, err := history.Load(histDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return err
	}

	results := history.DetectAnomalies(s, threshold)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Printf("No anomalies detected (threshold: %.1f%%)\n", threshold)
		return nil
	}

	for _, r := range results {
		fmt.Printf("Host: %s — rare ports (< %.1f%% of scans): %v\n", r.Host, r.ThresholdPct, r.RarePorts)
	}
	return nil
}
