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
		Use:   "pressure",
		Short: "Show port-pressure scores per host",
		Long: `Pressure measures the average number of open ports relative to the
peak observed, giving a normalised 0-1 score of how consistently
"loaded" each host's port surface is over time.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPressure(dir, minScans, format)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	cmd.Flags().IntVar(&minScans, "min-scans", 2, "minimum scans required to include a host")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runPressure(dir string, minScans int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Pressure(store, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no data")
		return nil
	}

	fmt.Printf("%-20s %10s %10s %10s %10s\n", "HOST", "AVG_PORTS", "PEAK", "PRESSURE", "SCANS")
	for _, r := range results {
		fmt.Printf("%-20s %10.2f %10d %10.3f %10d\n",
			r.Host, r.AvgOpenPorts, r.PeakPorts, r.Pressure, r.ScanCount)
	}
	return nil
}
