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
	var minRuns int
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "sequence",
		Short: "Show ports that appeared in consecutive scans",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSequence(dir, minRuns, jsonOut)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().IntVar(&minRuns, "min-runs", 3, "minimum consecutive appearances")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.AddCommand(cmd)
}

func runSequence(dir string, minRuns int, jsonOut bool) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Sequence(store, minRuns)

	if jsonOut {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no consecutive port sequences found")
		return nil
	}

	fmt.Printf("%-20s %6s %5s  %-25s %-25s\n", "HOST", "PORT", "RUNS", "FIRST SEEN", "LAST SEEN")
	for _, r := range results {
		fmt.Printf("%-20s %6d %5d  %-25s %-25s\n", r.Host, r.Port, r.Runs, r.Start, r.End)
	}
	return nil
}
