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
	var format string

	cmd := &cobra.Command{
		Use:   "rhythm",
		Short: "Analyse scan timing regularity per host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRhythm(dir, minScans, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "minimum scans required")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text|json")
	rootCmd.AddCommand(cmd)
}

func runRhythm(dir string, minScans int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}
	results := history.Rhythm(store, minScans)
	if len(results) == 0 {
		fmt.Println("no rhythm data available")
		return nil
	}
	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}
	fmt.Printf("%-20s %8s %12s %12s %s\n", "HOST", "SCANS", "AVG_SEC", "STDDEV_SEC", "REGULAR")
	for _, r := range results {
		reg := "no"
		if r.Regular {
			reg = "yes"
		}
		fmt.Printf("%-20s %8d %12.1f %12.1f %s\n",
			r.Host, r.ScanCount, r.AvgIntervalSec, r.StdDevSec, reg)
	}
	return nil
}
