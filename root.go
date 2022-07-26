package main

import "github.com/spf13/cobra"

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asciiedit",
		Short: "Edit asciicast files (produced by asciinema)",
	}
	cmd.AddCommand(speedupCmd())
	return cmd
}
