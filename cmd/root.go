package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "apiterm",
    Short: "A CLI tool to test APIs from the terminal",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Welcome to apiterm, Enter help for available commands.")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
