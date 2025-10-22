package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd — основная команда
var rootCmd = &cobra.Command{
	Use:   "cerber [domain]",
	Short: "Краткое описание",
	Long:  `Полное описание моего приложения`,
	Args:  cobra.ExactArgs(1), // Ожидаем ровно один аргумент
	Run: func(cmd *cobra.Command, args []string) {

		lookupHostRun(cmd, args)

		// Например, вызвать аналог find или lookup:
		// processDomain(domain)
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(findCmd)
	findCmd.AddCommand(findPathCmd)

	rootCmd.AddCommand(LookCmd)
}

// Execute запускает root команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cleanDomain(searchDomain string) string {
	if len(searchDomain) == 0 {
		panic("not domain")
	}
	// Убираем "http://", "https://", "www."
	searchDomain = strings.TrimPrefix(searchDomain, "http://")
	searchDomain = strings.TrimPrefix(searchDomain, "https://")
	searchDomain = strings.TrimSuffix(searchDomain, "/")

	if res := strings.HasPrefix(searchDomain, "www."); res {
		searchDomain = searchDomain[4:]
	}
	return searchDomain
}
