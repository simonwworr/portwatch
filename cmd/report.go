package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

var (
	reportFormat string
	reportOutput string
)

func init() {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Scan hosts and print a formatted change report",
		RunE:  runReport,
	}
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "text", "Output format: text or json")
	reportCmd.Flags().StringVarP(&reportOutput, "state", "s", "portwatch.state.json", "Path to state file")
	rootCmd.AddCommand(reportCmd)
}

func runReport(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("portwatch.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	previous, _ := state.Load(reportOutput)
	current := make(state.PortMap)

	for _, host := range cfg.Hosts {
		ports, err := scanner.OpenPorts(host, cfg.Ports, cfg.Timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: scan %s failed: %v\n", host, err)
			continue
		}
		current[host] = ports
	}

	if err := state.Save(reportOutput, current); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	diffs := make(map[string]state.Diff)
	for host, ports := range current {
		d := state.Diff(previous[host], ports)
		if len(d.Opened) > 0 || len(d.Closed) > 0 {
			diffs[host] = d
		}
	}

	fmt := report.Format(reportFormat)
	r := report.New(diffs)
	r.Print(os.Stdout, fmt)
	return nil
}
