package cmd

import (
	"fmt"
	"os"

	"github.com/SimiyuWafulah/apiterm/pkg/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apiterm",
	Short: "A CLI tool to test APIs from the terminal",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Launching APITERM TUI...")
		if err := tui.Run(); err != nil {
			fmt.Printf("Error launching TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}