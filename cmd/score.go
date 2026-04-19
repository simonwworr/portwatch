package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"portwatch/internal/history"
)

func init() {
	var format string
	var dir string

	cmd := &cobra.Command{
		Use:   "score",
		Short: "Show risk scores for monitored hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScore(dir, format)
		},
	}
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	cmd.Flags().StringVar(&dir, "dir", ".", "Directory containing watch log")
	rootCmd.AddCommand(cmd)
}

func runScore(dir, format string) error {
	scores, err := history.Score(dir)
	if err != nil {
		return fmt.Errorf("score: %w", err)
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(scores)
	}

	if len(scores) == 0 {
		fmt.Println("No data available.")
		return nil
	}

	fmt.Printf("%-20s %8s %8s %8s %8s\n", "HOST", "SCORE", "CHANGES", "OPENED", "CLOSED")
	for _, s := range scores {
		fmt.Printf("%-20s %8.2f %8d %8d %8d\n", s.Host, s.Score, s.Changes, s.Openings, s.Closings)
	}
	return nil
}
