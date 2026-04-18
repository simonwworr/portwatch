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
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "trend <host>",
		Short: "Show port frequency trends for a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTrend(args[0], historyFile, outputJSON)
		},
	}

	cmd.Flags().StringVar(&historyFile, "history", "portwatch-history.json", "Path to history file")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	rootCmd.AddCommand(cmd)
}

func runTrend(host, historyFile string, outputJSON bool) error {
	store, err := history.Load(historyFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}
	if store == nil {
		store = &history.Store{}
	}

	ht := history.Trend(store, host)

	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(ht)
	}

	if len(ht.Trends) == 0 {
		fmt.Printf("No trend data for host %s\n", host)
		return nil
	}

	fmt.Printf("Port trends for %s (%d scans):\n", host, ht.Trends[0].TotalScans)
	fmt.Printf("%-8s %-10s %-10s %-22s %-22s\n", "PORT", "SEEN", "FREQUENCY", "FIRST SEEN", "LAST SEEN")
	for _, tr := range ht.Trends {
		fmt.Printf("%-8d %-10d %-10.2f %-22s %-22s\n",
			tr.Port, tr.SeenCount, tr.Frequency, tr.FirstSeen, tr.LastSeen)
	}
	return nil
}
