package cmd

import (
	"log"

	"github.com/SimiyuWafulah/apiterm/pkg/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "apiterm",
	Short: "Launch the interactive Terminal UI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(); err != nil {
			log.Fatalf("Error running UI: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
