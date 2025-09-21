package cmd

import (
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/spf13/cobra"
)

// Defines a new Cobra command 
var getCmd = &cobra.Command{
    Use:   "get [url]",
    Short: "Send get request to any URL",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        url := args[0]
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Request error:", err)
            return
        }
        defer resp.Body.Close()

        // Reads the entire response body
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading response:", err)
            return
        }

        fmt.Println("Status:", resp.Status)
        fmt.Println(string(body))
    },
}

// Registers get  with the root command for use in the cli
func init() {
    rootCmd.AddCommand(getCmd)
}
