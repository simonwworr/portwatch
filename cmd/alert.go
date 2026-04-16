package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

var (
	alertWebhook  string
	alertStateDir string
)

func init() {
	alertCmd := &cobra.Command{
		Use:   "alert [host]",
		Short: "Scan a host and alert if ports have changed since last run",
		Args:  cobra.ExactArgs(1),
		RunE:  runAlert,
	}
	alertCmd.Flags().StringVar(&alertWebhook, "webhook", "", "Webhook URL to POST alerts to")
	alertCmd.Flags().StringVar(&alertStateDir, "state-dir", ".portwatch", "Directory to store state files")
	rootCmd.AddCommand(alertCmd)
}

func runAlert(cmd *cobra.Command, args []string) error {
	host := args[0]
	ports, err := scanner.OpenPorts(host)
	if err != nil {
		return fmt.Errorf("scan %s: %w", host, err)
	}

	if err := os.MkdirAll(alertStateDir, 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	stateFile := fmt.Sprintf("%s/%s.json", alertStateDir, host)
	prev, err := state.Load(stateFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load state: %w", err)
	}

	diff := state.Diff(prev, ports)

	if alert.HasChanges(diff) {
		var a alert.Alerter
		if alertWebhook != "" {
			a = alert.NewWebhookAlerter(alertWebhook)
		} else {
			a = alert.NewLogAlerter()
		}
		if err := a.Notify(host, diff); err != nil {
			return fmt.Errorf("notify: %w", err)
		}
	}

	return state.Save(stateFile, ports)
}
