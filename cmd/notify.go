package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

func init() {
	var cfgFile string
	var host string
	var opened string
	var closed string

	cmd := &cobra.Command{
		Use:   "notify",
		Short: "Send a test notification for a given host and diff",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNotify(cfgFile, host, opened, closed)
		},
	}

	cmd.Flags().StringVarP(&cfgFile, "config", "c", "portwatch.yaml", "Config file")
	cmd.Flags().StringVar(&host, "host", "", "Target host")
	cmd.Flags().StringVar(&opened, "opened", "", "Comma-separated opened ports")
	cmd.Flags().StringVar(&closed, "closed", "", "Comma-separated closed ports")
	_ = cmd.MarkFlagRequired("host")

	rootCmd.AddCommand(cmd)
}

func parsePorts(s string) []int {
	var ports []int
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err == nil {
			ports = append(ports, n)
		}
	}
	return ports
}

func runNotify(cfgFile, host, opened, closed string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	diff := state.Diff{
		Opened: parsePorts(opened),
		Closed: parsePorts(closed),
	}

	d := notify.NewDispatcher(notify.NewStdoutChannel())
	if cfg.Webhook != "" {
		log.Printf("webhook configured but skipped in test notify")
	}

	return d.Dispatch(host, diff)
}
