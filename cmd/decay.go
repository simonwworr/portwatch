package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var dir string
	var format string
	var halfLifeDays int

	cmd := &cobra.Command{
		Use:   "decay",
		Short: "Show exponential decay scores for observed ports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDecay(dir, format, halfLifeDays)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text|json")
	cmd.Flags().IntVar(&halfLifeDays, "half-life", 7, "half-life in days")
	rootCmd.AddCommand(cmd)
}

func runDecay(dir, format string, halfLifeDays int) error {
	store, err := history.Load(dir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}

	halfLife := time.Duration(halfLifeDays) * 24 * time.Hour
	results := history.Decay(store, halfLife, time.Now())

	sort.Slice(results, func(i, j int) bool {
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Score > results[j].Score
	})

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no port decay data available")
		return nil
	}
	fmt.Printf("%-20s %-8s %-12s %s\n", "HOST", "PORT", "AGE(days)", "SCORE")
	for _, r := range results {
		fmt.Printf("%-20s %-8d %-12.1f %.4f\n", r.Host, r.Port, r.AgeDays, r.Score)
	}
	return nil
}
