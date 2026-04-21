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
	var minScore float64
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "similarity",
		Short: "Show Jaccard similarity between host port profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSimilarity(dir, minScore, jsonOut)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	cmd.Flags().Float64Var(&minScore, "min", 0.0, "minimum similarity score (0.0–1.0)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.AddCommand(cmd)
}

func runSimilarity(dir string, minScore float64, jsonOut bool) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Similarity(store, minScore)
	if len(results) == 0 {
		fmt.Println("no similar host pairs found")
		return nil
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	fmt.Printf("%-20s %-20s %10s  %s\n", "HOST-A", "HOST-B", "SIMILARITY", "SHARED")
	for _, r := range results {
		fmt.Printf("%-20s %-20s %10.3f  %v\n", r.HostA, r.HostB, r.Similarity, r.Shared)
	}
	return nil
}
