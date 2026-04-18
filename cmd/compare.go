package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var historyFile string
	var format string

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare the two most recent scans per host and show port deltas",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(historyFile, format)
		},
	}
	cmd.Flags().StringVar(&historyFile, "history", "history.json", "Path to history file")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runCompare(historyFile, format string) error {
	store, err := history.Load(historyFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading history: %w", err)
	}
	if store == nil {
		store = &history.Store{}
	}

	deltas := history.Compare(store)

	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(deltas)
	}

	if len(deltas) == 0 {
		fmt.Println("No hosts with enough history to compare.")
		return nil
	}
	for _, d := range deltas {
		fmt.Printf("Host: %s\n", d.Host)
		if len(d.Opened) > 0 {
			fmt.Printf("  Opened: %v\n", d.Opened)
		}
		if len(d.Closed) > 0 {
			fmt.Printf("  Closed: %v\n", d.Closed)
		}
		if len(d.Stable) > 0 {
			fmt.Printf("  Stable: %v\n", d.Stable)
		}
	}
	return nil
}
