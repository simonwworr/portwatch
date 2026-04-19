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
	var minCount int
	var format string

	cmd := &cobra.Command{
		Use:   "correlation",
		Short: "Show ports that frequently appear together",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCorrelation(dir, minCount, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().IntVar(&minCount, "min", 2, "minimum co-occurrence count")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runCorrelation(dir string, minCount int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Correlate(store, minCount)
	if len(results) == 0 {
		fmt.Println("no correlated port pairs found")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	for _, r := range results {
		fmt.Printf("host: %s\n", r.Host)
		for _, p := range r.Pairs {
			fmt.Printf("  ports %d <-> %d  (seen together %d times)\n", p.PortA, p.PortB, p.Count)
		}
	}
	return nil
}
