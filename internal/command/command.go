package command

import (
	"fmt"
	"os"
	"strings"

	"cerber/internal/i18n"

	"github.com/spf13/cobra"
)

// rootCmd — основная команда
var rootCmd = &cobra.Command{
	Use:   "cerber",
	Short: i18n.T("cmd_short_root"),
	Long:  i18n.T("cmd_long_root"),
}

var cliLang string

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVar(&cliLang, "lang", "auto", i18n.T("lang_flag_desc"))
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return i18n.SetLang(cliLang)
	}
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
		return "", fmt.Errorf(i18n.T("err_domain_required"))
	}
	searchDomain = strings.TrimSpace(searchDomain)
	// Убираем "http://", "https://", "www." без учета регистра
	searchDomain = trimPrefixFold(searchDomain, "http://")
	searchDomain = trimPrefixFold(searchDomain, "https://")
	searchDomain = strings.TrimSuffix(searchDomain, "/")

	searchDomain = trimPrefixFold(searchDomain, "www.")
	if searchDomain == "" {
		return "", fmt.Errorf(i18n.T("err_domain_empty_after_normalize"))
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
