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
	var minScans int
	var format string

	cmd := &cobra.Command{
		Use:   "entropy",
		Short: "Compute Shannon entropy of open-port distributions per host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEntropy(dir, minScans, format)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "minimum scans required to include a host")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runEntropy(dir string, minScans int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Entropy(store, minScans)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no data")
		return nil
	}

	fmt.Printf("%-30s %10s %8s %8s\n", "HOST", "ENTROPY", "SCANS", "UNIQUE")
	fmt.Printf("%-30s %10s %8s %8s\n", "----", "-------", "-----", "------")
	for _, r := range results {
		fmt.Printf("%-30s %10.3f %8d %8d\n", r.Host, r.Entropy, r.Scans, r.Unique)
	}
	return nil
}
