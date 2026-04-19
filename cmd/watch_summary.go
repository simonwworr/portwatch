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
	var format string

	cmd := &cobra.Command{
		Use:   "watch-summary",
		Short: "Summarize watch log events per host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatchSummary(dir, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing watch log")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runWatchSummary(dir, format string) error {
	log, err := history.LoadWatchLog(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no watch log found")
			return nil
		}
		return err
	}

	summaries := history.SummarizeWatchLog(log)
	if len(summaries) == 0 {
		fmt.Println("no events recorded")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(summaries)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tTOTAL\tOPENED\tCLOSED\tTOP PORTS")
	for _, s := range summaries {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%v\n",
			s.Host, s.TotalEvents, s.Opened, s.Closed, s.TopPorts)
	}
	return w.Flush()
}
