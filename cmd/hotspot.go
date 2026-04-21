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
	var minHosts int
	var format string

	cmd := &cobra.Command{
		Use:   "hotspot",
		Short: "Show ports consistently open across multiple hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHotspot(dir, minHosts, format)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	cmd.Flags().IntVar(&minHosts, "min-hosts", 2, "minimum number of hosts a port must appear on")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")

	rootCmd.AddCommand(cmd)
}

func runHotspot(dir string, minHosts int, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	results := history.Hotspot(store, minHosts)
	if len(results) == 0 {
		fmt.Println("no hotspots found")
		return nil
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	fmt.Printf("%-8s %-12s %-10s %s\n", "PORT", "HOST_COUNT", "FREQUENCY", "HOSTS")
	for _, r := range results {
		hosts := ""
		for i, h := range r.Hosts {
			if i > 0 {
				hosts += ", "
			}
			hosts += h
		}
		fmt.Printf("%-8d %-12d %-10.2f %s\n", r.Port, r.HostCount, r.Frequency, hosts)
	}
	return nil
}
