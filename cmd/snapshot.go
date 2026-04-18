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
	var host string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Show the latest port snapshot for one or all hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshot(historyFile, host, jsonOut)
		},
	}

	cmd.Flags().StringVar(&historyFile, "history", "portwatch-history.json", "History file path")
	cmd.Flags().StringVar(&host, "host", "", "Filter to a specific host")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	rootCmd.AddCommand(cmd)
}

func runSnapshot(historyFile, host string, jsonOut bool) error {
	store, err := history.Load(historyFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}

	var snaps []history.Snapshot
	if host != "" {
		if s, ok := history.LatestSnapshot(store, host); ok {
			snaps = append(snaps, s)
		}
	} else {
		snaps = history.AllSnapshots(store)
	}

	if len(snaps) == 0 {
		fmt.Println("no snapshots found")
		return nil
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(snaps)
	}

	for _, s := range snaps {
		fmt.Printf("host: %s  time: %s  ports: %v\n", s.Host, s.Timestamp.Format("2006-01-02 15:04:05"), s.Ports)
	}
	return nil
}
