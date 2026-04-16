package cmd

import (
	"fmt"
	"time"

	"github.com/portwatch/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	scanStart   int
	scanEnd     int
	scanWorkers int
	scanTimeout int
)

var scanCmd = &cobra.Command{
	Use:   "scan [host]",
	Short: "Scan open TCP ports on a host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		if scanStart < 1 || scanEnd > 65535 || scanStart > scanEnd {
			return fmt.Errorf("invalid port range: %d-%d", scanStart, scanEnd)
		}

		timeout := time.Duration(scanTimeout) * time.Millisecond
		fmt.Printf("Scanning %s (ports %d-%d)...\n", host, scanStart, scanEnd)

		results := scanner.ScanHost(host, scanStart, scanEnd, scanWorkers, timeout)

		if len(results) == 0 {
			fmt.Println("No open ports found.")
			return nil
		}

		fmt.Printf("Open ports on %s:\n", host)
		for _, r := range results {
			fmt.Printf("  %d/tcp\n", r.Port)
		}
		return nil
	},
}

func init() {
	scanCmd.Flags().IntVar(&scanStart, "start", 1, "Start of port range")
	scanCmd.Flags().IntVar(&scanEnd, "end", 1024, "End of port range")
	scanCmd.Flags().IntVar(&scanWorkers, "workers", 100, "Number of concurrent workers")
	scanCmd.Flags().IntVar(&scanTimeout, "timeout", 500, "Connection timeout in milliseconds")
	rootCmd.AddCommand(scanCmd)
}
