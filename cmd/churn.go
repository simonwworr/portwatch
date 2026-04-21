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
		Use:   "churn",
		Short: "Show port churn rate per host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChurn(dir, minScans, format)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	cmd.Flags().IntVar(&minScans, "min-scans", 2, "minimum scans required to compute churn")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runChurn(dir string, minScans int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Churn(store, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no churn data available")
		return nil
	}

	fmt.Printf("%-20s %8s %8s %8s %10s\n", "HOST", "OPENED", "CLOSED", "TOTAL", "RATE")
	for _, r := range results {
		fmt.Printf("%-20s %8d %8d %8d %10.2f\n",
			r.Host, r.Opened, r.Closed, r.Total, r.ChurnRate)
	}
	return nil
}
