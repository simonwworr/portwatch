package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var (
		histDir  string
		minScans int
		format   string
	)

	cmd := &cobra.Command{
		Use:   "reachability",
		Short: "Show per-host uptime and port availability metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReachability(histDir, minScans, format)
		},
	}

	cmd.Flags().StringVar(&histDir, "dir", "history", "Directory containing history files")
	cmd.Flags().IntVar(&minScans, "min-scans", 2, "Minimum scans required to include a host")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runReachability(histDir string, minScans int, format string) error {
	store, err := history.Load(histDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return err
	}

	results := history.Reachability(store, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("No reachability data available.")
		return nil
	}

	fmt.Printf("%-20s %8s %8s %8s %8s %8s\n",
		"HOST", "SCANS", "UPTIME%", "AVG_PRT", "MAX_PRT", "MIN_PRT")
	fmt.Println("------------------------------------------------------------------------")
	for _, r := range results {
		fmt.Printf("%-20s %8d %7.1f%% %8.1f %8d %8d\n",
			r.Host, r.TotalScans, r.Uptime*100, r.AvgPorts, r.MaxPorts, r.MinPorts)
	}
	return nil
}
