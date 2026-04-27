package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

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
		Use:   "maturity",
		Short: "Score hosts by how mature (stable + long-lived) their port profile is",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMaturity(histDir, minScans, outFormat)
		},
	}

	cmd.Flags().StringVar(&histDir, "history-dir", "history", "Directory containing history JSON files")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "Minimum scans required to include a host")
	cmd.Flags().StringVar(&outFormat, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runMaturity(histDir string, minScans int, format string) error {
	store, err := history.Load(histDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return err
	}

	results := history.Maturity(store, minScans, time.Now())

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No maturity data available.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tSCORE\tAVG PORT AGE (days)\tSTABLE PORTS\tSCANS")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%.3f\t%.1f\t%v\t%d\n",
			r.Host, r.Score, r.AvgPortAge, r.StablePorts, r.ScanCount)
	}
	return w.Flush()
}
