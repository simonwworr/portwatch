package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/portwatch/internal/history"
)

func init() {
	var dir string
	var minHosts int
	var format string

	cmd := &cobra.Command{
		Use:   "cohort",
		Short: "Group hosts by the month they were first observed and show shared ports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCohort(dir, minHosts, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "history directory")
	cmd.Flags().IntVar(&minHosts, "min-hosts", 1, "minimum hosts per cohort")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runCohort(dir string, minHosts int, format string) error {
	store, err := history.Load(dir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}
	results := history.Cohort(store, minHosts)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no cohorts found")
		return nil
	}
	for _, c := range results {
		ports := make([]string, len(c.Ports))
		for i, p := range c.Ports {
			ports[i] = fmt.Sprintf("%d", p)
		}
		fmt.Printf("cohort %s  hosts=%s  shared_ports=[%s]  first_seen=%s\n",
			c.CohortID,
			strings.Join(c.Hosts, ","),
			strings.Join(ports, ","),
			c.FirstSeen.Format("2006-01-02"),
		)
	}
	return nil
}
