package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portwatch/internal/config"
)

func init() {
	var output string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a default portwatch config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(output)
		},
	}

	initCmd.Flags().StringVarP(&output, "output", "o", "portwatch.json", "Path to write the config file")
	rootCmd.AddCommand(initCmd)
}

func runInit(output string) error {
	if _, err := os.Stat(output); err == nil {
		return fmt.Errorf("config file %q already exists; remove it first", output)
	}

	cfg := config.DefaultConfig()
	cfg.Hosts = []string{"localhost"}
	cfg.Ports = []int{22, 80, 443, 8080}

	if err := config.Save(output, cfg); err != nil {
		return err
	}

	fmt.Printf("Created config file: %s\n", output)
	fmt.Println("Edit it to add your hosts and desired settings.")
	return nil
}
