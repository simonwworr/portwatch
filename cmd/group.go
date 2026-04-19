package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var bucket string
	var format string
	var dir string

	groupCmd := &cobra.Command{
		Use:   "group <host>",
		Short: "Group scan history by day, week, or month",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroup(args[0], bucket, format, dir)
		},
	}

	groupCmd.Flags().StringVarP(&bucket, "bucket", "b", "day", "Time bucket: day, week, month")
	groupCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	groupCmd.Flags().StringVarP(&dir, "dir", "d", ".", "History directory")

	rootCmd.AddCommand(groupCmd)
}

func runGroup(host, bucket, format, dir string) error {
	store, err := history.Load(dir)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	results, err := history.GroupBy(store, host, bucket)
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintf(os.Stdout, "no history for host %s\n", host)
		return nil
	}

	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	fmt.Fprintf(os.Stdout, "%-15s  %-6s  %-10s  %s\n", "bucket", "scans", "avg_ports", "unique_ports")
	for _, g := range results {
		fmt.Fprintf(os.Stdout, "%-15s  %-6d  %-10.1f  %v\n", g.Bucket, g.ScanCount, g.AvgPorts, g.UniquePorts)
	}
	return nil
}
