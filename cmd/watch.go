package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

var (
	watchInterval int
	watchStatePath string
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch [host]",
		Short: "Continuously monitor a host for port changes",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatch,
	}
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 60, "Scan interval in seconds")
	watchCmd.Flags().StringVarP(&watchStatePath, "state", "s", ".portwatch_state.json", "Path to state file")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	host := args[0]
	fmt.Printf("Watching %s every %ds. State: %s\n", host, watchInterval, watchStatePath)

	for {
		store, err := state.Load(watchStatePath)
		if err != nil {
			return fmt.Errorf("loading state: %w", err)
		}

		current, err := scanner.OpenPorts(host)
		if err != nil {
			fmt.Printf("[WARN] scan error for %s: %v\n", host, err)
		} else {
			sort.Ints(current)
			prev := store[host].Ports
			added, removed := state.Diff(prev, current)

			if len(added) > 0 {
				fmt.Printf("[ALERT] %s — new ports opened: %v\n", host, added)
			}
			if len(removed) > 0 {
				fmt.Printf("[ALERT] %s — ports closed: %v\n", host, removed)
			}
			if len(added) == 0 && len(removed) == 0 {
				fmt.Printf("[OK] %s — no changes\n", host)
			}

			store[host] = state.HostState{
				Host:      host,
				Ports:     current,
				ScannedAt: time.Now().UTC(),
			}
			if err := state.Save(watchStatePath, store); err != nil {
				return fmt.Errorf("saving state: %w", err)
			}
		}

		time.Sleep(time.Duration(watchInterval) * time.Second)
	}
}
