package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var days int
	var format string
	var dir string

	cmd := &cobra.Command{
		Use:   "pulse",
		Short: "Show activity pulse score per host",
		Long:  "Measures how consistently each host has been active over the observation window.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPulse(dir, days, format)
		},
	}

	cmd.Flags().IntVar(&days, "days", 30, "Observation window in days")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	cmd.Flags().StringVar(&dir, "dir", ".", "Directory containing history files")

	rootCmd.AddCommand(cmd)
}

func runPulse(dir string, days int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return fmt.Errorf("loading history: %w", err)
	}

	since := time.Now().AddDate(0, 0, -days)
	results := history.Pulse(store, since)

	if len(results) == 0 {
		fmt.Println("No pulse data available.")
		return nil
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(results)
	default:
		fmt.Printf("%-20s  %6s  %6s  %6s  %6s  %8s\n",
			"HOST", "PULSE", "SCANS", "ACTIVE", "TOTAL", "AVG_PORTS")
		fmt.Println("------------------------------------------------------------------")
		for _, r := range results {
			fmt.Printf("%-20s  %6.3f  %6d  %6d  %6d  %8.2f\n",
				r.Host, r.PulseScore, r.ScanCount, r.ActiveDays, r.TotalDays, r.AvgPorts)
		}
	}
	return nil
}
