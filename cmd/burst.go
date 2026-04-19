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
	var window int
	var threshold int
	var format string

	cmd := &cobra.Command{
		Use:   "burst <host>",
		Short: "Detect port-open bursts for a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBurst(args[0], histDir, window, threshold, format)
		},
	}

	cmd.Flags().StringVar(&histDir, "history-dir", "history", "Directory for history files")
	cmd.Flags().IntVar(&window, "window", 3, "Sliding window size (number of scans)")
	cmd.Flags().IntVar(&threshold, "threshold", 3, "Minimum new ports to count as a burst")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runBurst(host, histDir string, window, threshold int, format string) error {
	s, err := history.Load(histDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}

	results := history.DetectBursts(s, host, window, threshold)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Printf("No bursts detected for %s\n", host)
		return nil
	}
	for _, r := range results {
		fmt.Printf("host=%s burst_size=%d window=%s to %s ports=%v\n",
			r.Host, r.BurstSize,
			r.WindowStart.Format("2006-01-02 15:04"),
			r.WindowEnd.Format("2006-01-02 15:04"),
			r.PeakPorts)
	}
	return nil
}
