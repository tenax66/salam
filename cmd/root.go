package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "salam [OPTIONS] <URL>",
	Short: "Run benchmarks for <URL>",
	Long:  `Runs provided number of requests to <URL>.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		url := args[0]
		n, _ := cmd.Flags().GetInt("number")

		results := make(chan string, n)

		for i := 0; i < n; i++ {
			wg.Add(1)
			go fetchURL(&wg, url, results)
		}

		wg.Wait()
		close(results)

		for result := range results {
			fmt.Println(result)
		}
	},
}

func fetchURL(wg *sync.WaitGroup, url string, results chan<- string) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return
	}
	results <- fmt.Sprintf("Status code: %d, Time: %v", resp.StatusCode, duration)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.DisableFlagsInUseLine = true
	rootCmd.Flags().IntP("number", "n", 1, "number of requests to run")

	const usageTemplate = `Usage:
{{if .Runnable}}{{.UseLine}}{{end}}
{{if .HasAvailableSubCommands}}{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
	rootCmd.SetUsageTemplate(usageTemplate)
}
