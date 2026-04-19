package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var format string
	var threshold float64
	var dir string

	cmd := &cobra.Command{
		Use:   "spike <host>",
		Short: "Detect port count spikes for a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSpike(args[0], dir, threshold, format)
		},
	}
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	cmd.Flags().Float64Var(&threshold, "threshold", 2.0, "Spike threshold multiplier")
	cmd.Flags().StringVar(&dir, "dir", ".", "History directory")
	rootCmd.AddCommand(cmd)
}

func runSpike(host, dir string, threshold float64, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.DetectSpikes(store, host, threshold)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Printf("no spikes detected for %s (threshold %.1fx)\n", host, threshold)
		return nil
	}
	fmt.Printf("spikes for %s (threshold %.1fx):\n", host, threshold)
	for _, r := range results {
		fmt.Printf("  %s  ports=%d  avg=%.1f  delta=+%.1f\n",
			r.At.Format("2006-01-02 15:04"), r.PortCount, r.Avg, r.Delta)
	}
	return nil
}
