package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portwatch",
	Short: "Lightweight to monitor and alert on open port changes acrossong: `ans TCP ports on one or more hosts and alerts younwhen the set previous snapshot.`,
}unc Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
