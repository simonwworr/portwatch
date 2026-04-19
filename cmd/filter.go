package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var host string
	var port int
	var since, until string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter history entries by host, port, or time range",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			store, err := history.Load(dir)
			if err != nil {
				return fmt.Errorf("load history: %w", err)
			}
			opts := history.FilterOptions{Host: host, Port: port}
			if since != "" {
				if t, err := time.Parse(time.RFC3339, since); err == nil {
					opts.Since = t
				}
			}
			if until != "" {
				if t, err := time.Parse(time.RFC3339, until); err == nil {
					opts.Until = t
				}
			}
			entries := history.Filter(store, opts)
			if jsonOut {
				return json.NewEncoder(os.Stdout).Encode(entries)
			}
			if len(entries) == 0 {
				fmt.Println("no entries matched")
				return nil
			}
			for _, e := range entries {
				fmt.Printf("%s  %s  ports=%v\n", e.Time.Format(time.RFC3339), e.Host, e.Ports)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "filter by host")
	cmd.Flags().IntVar(&port, "port", 0, "filter by port")
	cmd.Flags().StringVar(&since, "since", "", "filter entries after this time (RFC3339)")
	cmd.Flags().StringVar(&until, "until", "", "filter entries before this time (RFC3339)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")

	rootCmd.AddCommand(cmd)
}
