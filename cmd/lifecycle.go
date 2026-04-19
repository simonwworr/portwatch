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
	var host string

	cmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Show port open/close lifecycle spans for a host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLifecycle(dir, host, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().StringVar(&host, "host", "", "host to inspect (required)")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text|json")
	_ = cmd.MarkFlagRequired("host")
	rootCmd.AddCommand(cmd)
}

func runLifecycle(dir, host, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	events := history.Lifecycle(store, host)
	if len(events) == 0 {
		fmt.Printf("no lifecycle data for %s\n", host)
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(events)
	}

	fmt.Printf("%-8s  %-24s  %-24s  %s\n", "PORT", "OPENED", "CLOSED", "DURATION")
	for _, e := range events {
		closed := "(open)"
		dur := "-"
		if e.ClosedAt != nil {
			closed = e.ClosedAt.Format("2006-01-02 15:04:05")
			dur = e.Duration.Round(1e9).String()
		}
		fmt.Printf("%-8d  %-24s  %-24s  %s\n",
			e.Port,
			e.OpenedAt.Format("2006-01-02 15:04:05"),
			closed,
			dur,
		)
	}
	return nil
}
