package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"portwatch/internal/history"
)

var pinDir string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin ports for a host with an optional note",
	}

	addCmd := &cobra.Command{
		Use:   "add <host> <port,...> [note]",
		Short: "Pin one or more ports for a host",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  runPinAdd,
	}

	listCmd := &cobra.Command{
		Use:   "list [host]",
		Short: "List pinned ports",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runPinList,
	}

	rmCmd := &cobra.Command{
		Use:   "remove <host> <port>",
		Short: "Remove a pinned port",
		Args:  cobra.ExactArgs(2),
		RunE:  runPinRemove,
	}

	pinCmd.PersistentFlags().StringVar(&pinDir, "dir", ".portwatch", "directory for pin storage")
	pinCmd.AddCommand(addCmd, listCmd, rmCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinAdd(cmd *cobra.Command, args []string) error {
	host := args[0]
	var ports []int
	for _, p := range strings.Split(args[1], ",") {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("invalid port %q: %w", p, err)
		}
		ports = append(ports, n)
	}
	note := ""
	if len(args) == 3 {
		note = args[2]
	}
	s := history.NewPinStore(pinDir)
	_ = s.Load()
	s.Add(host, ports, note)
	if err := s.Save(); err != nil {
		return err
	}
	fmt.Printf("Pinned %v for %s\n", ports, host)
	return nil
}

func runPinList(cmd *cobra.Command, args []string) error {
	s := history.NewPinStore(pinDir)
	if err := s.Load(); err != nil {
		return err
	}
	var pins []history.Pin
	if len(args) == 1 {
		pins = s.ForHost(args[0])
	} else {
		pins = s.All()
	}
	if len(pins) == 0 {
		fmt.Println("No pinned ports.")
		return nil
	}
	for _, p := range pins {
		fmt.Printf("%-20s ports=%-20v note=%s\n", p.Host, p.Ports, p.Note)
	}
	return nil
}

func runPinRemove(cmd *cobra.Command, args []string) error {
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	s := history.NewPinStore(pinDir)
	if err := s.Load(); err != nil {
		return err
	}
	if !s.Remove(args[0], port) {
		fmt.Println("No matching pin found.")
		return nil
	}
	if err := s.Save(); err != nil {
		return err
	}
	fmt.Printf("Removed pin for port %d on %s\n", port, args[0])
	return nil
}
