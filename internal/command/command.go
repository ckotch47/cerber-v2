package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd — основная команда
var rootCmd = &cobra.Command{
	Use:   "cerber",
	Short: "CLI для DNS recon и поиска скрытых путей",
	Long:  `Cerber — инструмент для DNS lookup, поиска поддоменов и скрытых путей.`,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(findCmd)
	findCmd.AddCommand(findPathCmd)

	rootCmd.AddCommand(LookCmd)
	rootCmd.AddCommand(googleCmd)
	rootCmd.AddCommand(apiCmd)
}

// Execute запускает root команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cleanDomain(searchDomain string) (string, error) {
	if len(strings.TrimSpace(searchDomain)) == 0 {
		return "", fmt.Errorf("домен не указан")
	}
	searchDomain = strings.TrimSpace(searchDomain)
	// Убираем "http://", "https://", "www." без учета регистра
	searchDomain = trimPrefixFold(searchDomain, "http://")
	searchDomain = trimPrefixFold(searchDomain, "https://")
	searchDomain = strings.TrimSuffix(searchDomain, "/")

	searchDomain = trimPrefixFold(searchDomain, "www.")
	if searchDomain == "" {
		return "", fmt.Errorf("домен пуст после нормализации")
	}
	return searchDomain, nil
}

func trimPrefixFold(s, prefix string) string {
	if len(s) < len(prefix) {
		return s
	}
	if strings.EqualFold(s[:len(prefix)], prefix) {
		return s[len(prefix):]
	}
	return s
}
