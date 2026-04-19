package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"portwatch/internal/history"
)

var bookmarkDir string

func init() {
	bookmarkCmd := &cobra.Command{
		Use:   "bookmark",
		Short: "Manage port bookmarks",
	}

	addCmd := &cobra.Command{
		Use:   "add <host> <port> <label>",
		Short: "Bookmark a port on a host",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %w", err)
			}
			bs, err := history.NewBookmarkStore(bookmarkDir)
			if err != nil {
				return err
			}
			bs.Add(args[0], port, args[2])
			return bs.Save()
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all bookmarks",
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := history.NewBookmarkStore(bookmarkDir)
			if err != nil {
				return err
			}
			for _, b := range bs.All() {
				fmt.Printf("%s\t%d\t%s\n", b.Host, b.Port, b.Label)
			}
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <host> <port>",
		Short: "Remove a bookmark",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %w", err)
			}
			bs, err := history.NewBookmarkStore(bookmarkDir)
			if err != nil {
				return err
			}
			bs.Remove(args[0], port)
			return bs.Save()
		},
	}

	bookmarkCmd.PersistentFlags().StringVar(&bookmarkDir, "dir", ".portwatch", "directory for bookmark storage")
	bookmarkCmd.AddCommand(addCmd, listCmd, removeCmd)
	rootCmd.AddCommand(bookmarkCmd)
}
