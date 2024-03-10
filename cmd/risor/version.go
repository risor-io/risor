package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Risor",
	Run: func(cmd *cobra.Command, args []string) {
		outFmt := cmd.Flag("output").Value.String()
		if strings.ToLower(outFmt) == "json" {
			info, err := json.MarshalIndent(map[string]interface{}{
				"version": version,
				"commit":  commit,
				"date":    date,
			}, "", "  ")
			if err != nil {
				fatal(err)
			}
			fmt.Println(string(info))
		} else {
			fmt.Println(version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringP("output", "o", "", "Set the output format")
	versionCmd.RegisterFlagCompletionFunc("output",
		cobra.FixedCompletions(
			outputFormatsCompletion,
			cobra.ShellCompDirectiveNoFileComp,
		))
}
