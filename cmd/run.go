package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [OPTIONS] <URL>",
	Short: "Run benchmarks for <URL>.",
	Long:  `Runs provided number of requests to <URL>.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		n, _ := cmd.Flags().GetInt("number")
		fmt.Printf("%v times %v\n", url, n)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.DisableFlagsInUseLine = true
	runCmd.Flags().IntP("number", "n", 1, "number of requests to run")
}
