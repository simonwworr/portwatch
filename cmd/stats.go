package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var stateFile string
	var topN int

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show port frequency statistics from scan history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(stateFile, topN)
		},
	}
	cmd.Flags().StringVar(&stateFile, "history", "portwatch-history.json", "Path to history file")
	cmd.Flags().IntVar(&topN, "top", 5, "Number of top ports to display per host (0 = all)")
	rootCmd.AddCommand(cmd)
}

func runStats(historyFile string, topN int) error {
	s, err := history.Load(historyFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading history: %w", err)
	}
	if s == nil {
		fmt.Println("No history found.")
		return nil
	}

	results := history.Stats(s, topN)
	if len(results) == 0 {
		fmt.Println("No data available.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tSCANS\tTOP PORTS (seen count)")
	for _, hs := range results {
		ports := ""
		for i, ps := range hs.TopPorts {
			if i > 0 {
				ports += ", "
			}
			ports += fmt.Sprintf("%d(x%d)", ps.Port, ps.SeenCount)
		}
		if ports == "" {
			ports = "-"
		}
		fmt.Fprintf(w, "%s\t%d\t%s\n", hs.Host, hs.ScanCount, ports)
	}
	return w.Flush()
}
