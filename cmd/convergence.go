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
	var format string

	cmd := &cobra.Command{
		Use:   "convergence",
		Short: "Report how quickly each host's port set stabilises",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConvergence(dir, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text|json")
	rootCmd.AddCommand(cmd)
}

func runConvergence(dir, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Convergence(store)
	if len(results) == 0 {
		fmt.Println("no convergence data (need at least 2 scans per host)")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	fmt.Printf("%-20s %6s %12s %7s %10s\n", "HOST", "SCANS", "STABLE_AFTER", "STABLE", "CHNG_RATE")
	for _, r := range results {
		stable := "no"
		if r.Stable {
			stable = "yes"
		}
		fmt.Printf("%-20s %6d %12d %7s %9.1f%%\n",
			r.Host, r.Scans, r.StableAfter, stable, r.ChangeRate*100)
	}
	return nil
}
