package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
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
		for i := 0; i < n; i++ {
			start := time.Now()
			status, _, err := fasthttp.Get(nil, url)
			duration := time.Since(start)

			if err != nil {
				return
			}
			fmt.Printf("Status code: %d, Time: %v", status, duration)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.DisableFlagsInUseLine = true
	runCmd.Flags().IntP("number", "n", 1, "number of requests to run")
}
