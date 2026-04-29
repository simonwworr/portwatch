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
		histDir   string
		minScans  int
		outFormat string
	)

	cmd := &cobra.Command{
		Use:   "flux",
		Short: "Show per-host port flux (rate of change between scans)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFlux(histDir, minScans, outFormat)
		},
	}

	cmd.Flags().StringVar(&histDir, "history-dir", "history", "Directory containing history JSON files")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "Minimum scans required to include a host")
	cmd.Flags().StringVar(&outFormat, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runFlux(histDir string, minScans int, format string) error {
	st, err := history.Load(histDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return err
	}

	results := history.Flux(st, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No flux data available (insufficient scans).")
		return nil
	}

	fmt.Printf("%-20s  %8s  %10s  %6s\n", "HOST", "FLUX", "AVG DELTA", "SCANS")
	fmt.Printf("%-20s  %8s  %10s  %6s\n", "----", "----", "---------", "-----")
	for _, r := range results {
		fmt.Printf("%-20s  %8.3f  %10.2f  %6d\n", r.Host, r.Flux, r.AvgDelta, r.Scans)
	}
	return nil
}
