package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"portwatch/internal/history"
)

func init() {
	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage port notes",
	}

	addCmd := &cobra.Command{
		Use:   "add <host> <port> <text>",
		Short: "Add a note to a host/port",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %w", err)
			}
			dir, _ := cmd.Flags().GetString("dir")
			ns, err := history.NewNoteStore(dir)
			if err != nil {
				return err
			}
			if err := ns.Add(host, port, args[2]); err != nil {
				return err
			}
			fmt.Printf("Note added for %s:%d\n", host, port)
			return nil
		},
	}
	addCmd.Flags().String("dir", ".portwatch", "history directory")

	listCmd := &cobra.Command{
		Use:   "list <host>",
		Short: "List notes for a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			ns, err := history.NewNoteStore(dir)
			if err != nil {
				return err
			}
			notes := ns.ForHost(args[0])
			if len(notes) == 0 {
				fmt.Println("No notes found.")
				return nil
			}
			for _, n := range notes {
				fmt.Printf("[%s] port %d: %s\n", n.CreatedAt.Format("2006-01-02 15:04"), n.Port, n.Text)
			}
			return nil
		},
	}
	listCmd.Flags().String("dir", ".portwatch", "history directory")

	noteCmd.AddCommand(addCmd, listCmd)
	rootCmd.AddCommand(noteCmd)
}
