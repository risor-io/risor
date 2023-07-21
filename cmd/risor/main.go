package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "dev"

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
			fmt.Println(version)
		},
	}

	rootCmd.AddCommand(cmdServe)
	rootCmd.AddCommand(cmdVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
