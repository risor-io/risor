package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Run the Risor API server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Server")
		},
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print the version of Risor",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			outFmt := cmd.Flag("output").Value.String()
			if strings.ToLower(outFmt) == "json" {
				info, err := json.MarshalIndent(map[string]interface{}{
					"version": version,
					"commit":  commit,
					"date":    date,
				}, "", "  ")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(string(info))
			} else {
				fmt.Println(version)
			}
		},
	}

	cmdVersion.Flags().StringP("output", "o", "", "Set the output format")

	rootCmd.AddCommand(cmdServe)
	rootCmd.AddCommand(cmdVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
