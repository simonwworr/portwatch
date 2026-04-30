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
	var jitterThreshold float64
	var format string
	var onlyIrregular bool

	cmd := &cobra.Command{
		Use:   "cadence",
		Short: "Analyse scan timing regularity per host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCadence(dir, minScans, jitterThreshold, format, onlyIrregular)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "history", "history directory")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "minimum scans required")
	cmd.Flags().Float64Var(&jitterThreshold, "jitter", 30.0, "jitter threshold in seconds for 'regular' classification")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	cmd.Flags().BoolVar(&onlyIrregular, "irregular", false, "show only irregular hosts")

	rootCmd.AddCommand(cmd)
}

func runCadence(dir string, minScans int, jitterThreshold float64, format string, onlyIrregular bool) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Cadence(store, minScans, jitterThreshold)

	if onlyIrregular {
		filtered := results[:0]
		for _, r := range results {
			if !r.Regular {
				filtered = append(filtered, r)
			}
		}
		results = filtered
	}

	if len(results) == 0 {
		fmt.Println("no cadence data")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	fmt.Printf("%-20s %6s %12s %12s %12s %8s %s\n",
		"HOST", "SCANS", "AVG", "MIN", "MAX", "JITTER", "REGULAR")
	for _, r := range results {
		reg := "yes"
		if !r.Regular {
			reg = "no"
		}
		fmt.Printf("%-20s %6d %12s %12s %12s %8.2f %s\n",
			r.Host, r.ScanCount,
			r.AvgInterval.Round(1),
			r.MinInterval.Round(1),
			r.MaxInterval.Round(1),
			r.Jitter, reg)
	}
	return nil
}
