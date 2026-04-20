package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var dir string
	var host string
	var format string

	cmd := &cobra.Command{
		Use:   "session",
		Short: "List recorded watch sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSession(dir, host, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".portwatch", "History directory")
	cmd.Flags().StringVar(&host, "host", "", "Filter by host")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runSession(dir, host, format string) error {
	all, err := history.LoadSessions(dir)
	if err != nil {
		return fmt.Errorf("load sessions: %w", err)
	}

	var sessions []history.Session
	if host != "" {
		sessions = history.SessionsForHost(all, host)
	} else {
		sessions = all
	}

	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(sessions)
	}

	if len(sessions) == 0 {
		fmt.Println("no sessions recorded")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tHOST\tSTARTED\tENDED\tSCANS\tPORTS")
	for _, s := range sessions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%v\n",
			s.ID,
			s.Host,
			s.StartedAt.Format("2006-01-02 15:04:05"),
			s.EndedAt.Format("2006-01-02 15:04:05"),
			s.ScanCount,
			s.PortsSeen,
		)
	}
	return w.Flush()
}
