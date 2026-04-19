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
		Use:   "heatmap",
		Short: "Show day-bucketed port activity heatmap",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHeatmap(dir, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runHeatmap(dir, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	entries := history.Heatmap(store)
	if len(entries) == 0 {
		fmt.Println("no heatmap data")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(entries)
	}

	fmt.Printf("%-20s %-12s %s\n", "HOST", "DATE", "PORTS")
	for _, e := range entries {
		fmt.Printf("%-20s %-12s %d\n", e.Host, e.Bucket, e.Count)
	}
	return nil
}
