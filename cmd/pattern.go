package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var historyDir string
	var minScans int
	var format string

	cmd := &cobra.Command{
		Use:   "pattern",
		Short: "Show port appearance frequency patterns across scan history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPattern(historyDir, minScans, format)
		},
	}
	cmd.Flags().StringVar(&historyDir, "history-dir", ".portwatch/history", "History directory")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "Minimum scans required to include a host")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runPattern(historyDir string, minScans int, format string) error {
	store, err := history.Load(historyDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}
	if store == nil {
		store = history.NewMemoryStore()
	}

	patterns := history.Pattern(store, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(patterns)
	}

	if len(patterns) == 0 {
		fmt.Println("No pattern data available.")
		return nil
	}

	fmt.Printf("%-20s %6s %6s %6s %8s\n", "HOST", "PORT", "SEEN", "TOTAL", "FREQ")
	for _, p := range patterns {
		fmt.Printf("%-20s %6d %6d %6d %7.1f%%\n",
			p.Host, p.Port, p.SeenCount, p.TotalScans, p.Frequency*100)
	}
	return nil
}
