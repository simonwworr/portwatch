package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var historyFile string
	var host string
	var pruneAge time.Duration
	var maxEntries int

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show or prune scan history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistory(historyFile, host, pruneAge, maxEntries)
		},
	}
	cmd.Flags().StringVar(&historyFile, "file", "portwatch-history.json", "Path to history file")
	cmd.Flags().StringVar(&host, "host", "", "Filter output to a specific host")
	cmd.Flags().DurationVar(&pruneAge, "prune-age", 0, "Remove entries older than this duration (e.g. 168h)")
	cmd.Flags().IntVar(&maxEntries, "prune-max", 0, "Keep at most N entries per host")
	rootCmd.AddCommand(cmd)
}

func runHistory(path, host string, pruneAge time.Duration, maxEntries int) error {
	l, err := history.Load(path)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	if pruneAge > 0 || maxEntries > 0 {
		n := history.Prune(l, history.PruneOptions{MaxAge: pruneAge, MaxEntries: maxEntries})
		if err := history.Save(path, l); err != nil {
			return fmt.Errorf("saving history: %w", err)
		}
		fmt.Printf("Pruned %d entries.\n", n)
		return nil
	}

	entries := l.Entries
	if host != "" {
		entries = l.ForHost(host)
	}
	if len(entries) == 0 {
		fmt.Println("No history entries found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tHOST\tPORTS")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%v\n", e.Timestamp.Format(time.RFC3339), e.Host, e.Ports)
	}
	return w.Flush()
}
