package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"portwatch/internal/history"
)

func init() {
	var dir string
	var format string

	cmd := &cobra.Command{
		Use:   "topology",
		Short: "Show port-sharing graph across scanned hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTopology(dir, format)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "directory containing history files")
	cmd.Flags().StringVar(&format, "format", "text", "output format: text or json")
	rootCmd.AddCommand(cmd)
}

func runTopology(dir, format string) error {
	store, err := history.Load(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no history found")
			return nil
		}
		return err
	}

	res := history.Topology(store)

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(res)
	}

	if len(res.Nodes) == 0 {
		fmt.Println("no hosts in history")
		return nil
	}

	fmt.Println("=== Hosts ===")
	for _, n := range res.Nodes {
		ports := make([]string, len(n.Ports))
		for i, p := range n.Ports {
			ports[i] = fmt.Sprintf("%d", p)
		}
		fmt.Printf("  %-20s  ports: [%s]\n", n.Host, strings.Join(ports, ", "))
	}

	if len(res.Edges) == 0 {
		fmt.Println("\nno shared ports between hosts")
		return nil
	}

	fmt.Println("\n=== Shared-Port Edges ===")
	for _, e := range res.Edges {
		ports := make([]string, len(e.SharedPorts))
		for i, p := range e.SharedPorts {
			ports[i] = fmt.Sprintf("%d", p)
		}
		fmt.Printf("  %s <-> %s  shared: [%s]\n", e.HostA, e.HostB, strings.Join(ports, ", "))
	}
	return nil
}
