package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var histDir string

	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage port labels",
	}

	addCmd := &cobra.Command{
		Use:   "add <host> <port> <label>",
		Short: "Add a label to a port",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %w", err)
			}
			s, err := history.LoadLabels(histDir)
			if err != nil {
				return err
			}
			if err := s.Add(args[0], port, args[2]); err != nil {
				return err
			}
			fmt.Printf("Labeled %s:%d as %q\n", args[0], port, args[2])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <host>",
		Short: "List labels for a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := history.LoadLabels(histDir)
			if err != nil {
				return err
			}
			labels := s.ForHost(args[0])
			if len(labels) == 0 {
				fmt.Printf("No labels for %s\n", args[0])
				return nil
			}
			for _, l := range labels {
				fmt.Printf("  :%d  %s  (%s)\n", l.Port, l.Label, l.CreatedAt.Format("2006-01-02"))
			}
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, listCmd} {
		sub.Flags().StringVar(&histDir, "history-dir", ".portwatch/history", "History directory")
		cmd.AddCommand(sub)
	}
	rootCmd.AddCommand(cmd)
}
