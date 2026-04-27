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
		dir      string
		minScans int
		lambda   float64
		format   string
	)

	cmd := &cobra.Command{
		Use:   "momentum",
		Short: "Show weighted port-change momentum per host",
		Long: `Computes a momentum score for each host based on the weighted rate of port
count change over time. Positive scores indicate ports are trending open;
negative scores indicate ports are closing. Recent changes are weighted
more heavily using exponential decay.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMomentum(dir, minScans, lambda, format)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "./history", "History directory")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "Minimum scans required")
	cmd.Flags().Float64Var(&lambda, "lambda", 0.5, "Exponential decay factor (higher = more weight on recent)")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runMomentum(dir string, minScans int, lambda float64, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return fmt.Errorf("loading history: %w", err)
	}

	results := history.Momentum(store, minScans, lambda)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No momentum data available (insufficient scans).")
		return nil
	}

	fmt.Printf("%-30s %10s %10s %6s\n", "HOST", "MOMENTUM", "NET CHANGE", "SCANS")
	fmt.Printf("%-30s %10s %10s %6s\n", "----", "--------", "----------", "-----")
	for _, r := range results {
		direction := " "
		if r.Score > 0 {
			direction = "↑"
		} else if r.Score < 0 {
			direction = "↓"
		}
		fmt.Printf("%-30s %9.3f%s %+10d %6d\n", r.Host, r.Score, direction, r.NetChange, r.Scans)
	}
	return nil
}
