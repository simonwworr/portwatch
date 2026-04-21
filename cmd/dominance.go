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
		minPorts int
		format   string
	)

	cmd := &cobra.Command{
		Use:   "dominance",
		Short: "Rank hosts by how often their port set is a superset of peers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDominance(histDir, minPorts, format)
		},
	}

	cmd.Flags().StringVar(&histDir, "history-dir", ".portwatch/history", "Directory containing history files")
	cmd.Flags().IntVar(&minPorts, "min-ports", 1, "Minimum open ports for a host to be evaluated")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runDominance(histDir string, minPorts int, format string) error {
	store, err := history.Load(histDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No history found.")
			return nil
		}
		return fmt.Errorf("loading history: %w", err)
	}

	results := history.Dominance(store, minPorts)
	if len(results) == 0 {
		fmt.Println("No dominance data available.")
		return nil
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(results)
	default:
		fmt.Printf("%-30s  %6s  %5s\n", "HOST", "SCORE", "PEERS")
		for _, r := range results {
			fmt.Printf("%-30s  %6.2f  %5d\n", r.Host, r.Score, r.Peers)
		}
	}
	return nil
}
