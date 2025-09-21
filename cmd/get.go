package cmd

import (
    "fmt"
    "github.com/SimiyuWafulah/apiterm/internal"
    "github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
    Use:   "get [url]",
    Short: "Send a GET request to the specified URL",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        url := args[0]
        status, body, err := internal.Get(url)
        if err != nil {
            fmt.Println("Request error:", err)
            return
        }
        fmt.Println("Status:", status)
        fmt.Println(string(body))
    },
}
