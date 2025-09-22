package cmd

import (
    _"bytes"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/SimiyuWafulah/apiterm/internal"
)

//Defines Post command
var postCmd = &cobra.Command{
    Use:   "post [url] [json]",
    Short: "Send a post request with a JSON body",
    Args:  cobra.MinimumNArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        url := args[0]
        jsonBody := args[1]
        status, body, err := internal.Post(url, []byte(jsonBody))
        if err != nil {
            fmt.Println("Request error:", err)
            return
        }
        fmt.Println("Status:", status)
        fmt.Println(string(body))
    },
}

func init() {
    rootCmd.AddCommand(postCmd)
}
