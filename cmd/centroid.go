package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var dir string
	var minScans int
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "centroid",
		Short: "Show the majority port centroid and deviation for each host",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCentroid(dir, minScans, jsonOut)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	cmd.Flags().IntVar(&minScans, "min-scans", 3, "minimum scans required to include a host")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")

	rootCmd.AddCommand(cmd)
}

func runCentroid(dir string, minScans int, jsonOut bool) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Centroid(store, minScans)

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no hosts meet the minimum scan threshold")
		return nil
	}

	for _, r := range results {
		ports := make([]string, len(r.Centroid))
		for i, p := range r.Centroid {
			ports[i] = fmt.Sprintf("%d", p)
		}
		centroidStr := strings.Join(ports, ", ")
		if centroidStr == "" {
			centroidStr = "(none)"
		}
		fmt.Printf("%-20s deviation=%.3f centroid=[%s]\n", r.Host, r.Deviation, centroidStr)
	}
	return nil
}
