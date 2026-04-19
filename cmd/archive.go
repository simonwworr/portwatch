package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"portwatch/internal/history"
)

func init() {
	var dir string
	var maxAgeDays int
	var host string
	var listOnly bool

	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive old history entries or list existing archives",
		RunE: func(cmd *cobra.Command, args []string) error {
			if listOnly {
				return runListArchives(dir, host)
			}
			return runArchive(dir, host, maxAgeDays)
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "History directory")
	cmd.Flags().IntVar(&maxAgeDays, "max-age", 30, "Archive entries older than this many days")
	cmd.Flags().StringVar(&host, "host", "", "Host to archive (required)")
	cmd.Flags().BoolVar(&listOnly, "list", false, "List existing archives instead of archiving")
	_ = cmd.MarkFlagRequired("host")

	rootCmd.AddCommand(cmd)
}

func runArchive(dir, host string, maxAgeDays int) error {
	store, err := history.Load(dir)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	maxAge := time.Duration(maxAgeDays) * 24 * time.Hour
	n, err := history.Archive(dir, store, host, maxAge)
	if err != nil {
		return fmt.Errorf("archive: %w", err)
	}

	if err := history.Save(dir, store); err != nil {
		return fmt.Errorf("save history: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Archived %d entries for %s (older than %d days)\n", n, host, maxAgeDays)
	return nil
}

func runListArchives(dir, host string) error {
	archives, err := history.LoadArchives(dir, host)
	if err != nil {
		return fmt.Errorf("load archives: %w", err)
	}
	if len(archives) == 0 {
		fmt.Fprintf(os.Stdout, "No archives found for %s\n", host)
		return nil
	}
	for _, a := range archives {
		fmt.Fprintf(os.Stdout, "[%s] host=%s entries=%d\n",
			a.ArchivedAt.Format(time.RFC3339), a.Host, len(a.Entries))
	}
	return nil
}
