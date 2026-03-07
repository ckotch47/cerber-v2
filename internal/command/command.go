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
		if err := i18n.SetLang(cliLang); err != nil {
			return err
		}
		applyLocalizedTexts()
		return nil
	}
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(findCmd)
	findCmd.AddCommand(findPathCmd)

	rootCmd.AddCommand(LookCmd)
	rootCmd.AddCommand(googleCmd)
	rootCmd.AddCommand(apiCmd)

	applyLocalizedTexts()
}

// Execute запускает root команду
func Execute() {
	if err := i18n.SetLang(langFromArgs(os.Args[1:])); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	applyLocalizedTexts()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func langFromArgs(args []string) string {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--lang" && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, "--lang=") {
			return strings.TrimPrefix(arg, "--lang=")
		}
	}
	return "auto"
}

func applyLocalizedTexts() {
	rootCmd.Short = i18n.T("cmd_short_root")
	rootCmd.Long = i18n.T("cmd_long_root")
	setFlagUsage(rootCmd, "lang", "lang_flag_desc")

	findCmd.Short = i18n.T("cmd_short_find")
	setFlagUsage(findCmd, "wordlist", "flag_wordlist_desc")
	setFlagUsage(findCmd, "worldlis", "flag_worldlis_desc")
	setFlagUsage(findCmd, "recurse", "flag_recurse_desc")
	setFlagUsage(findCmd, "max-depth", "flag_max_depth_desc")
	setFlagUsage(findCmd, "concurrency", "flag_concurrency_desc")

	findPathCmd.Short = i18n.T("cmd_short_find_path")
	findPathCmd.Long = i18n.T("cmd_short_find_path")
	setFlagUsage(findPathCmd, "wordlist", "flag_wordlist_desc")
	setFlagUsage(findPathCmd, "worldlis", "flag_worldlis_desc")
	setFlagUsage(findPathCmd, "exclude", "flag_exclude_desc")
	setFlagUsage(findPathCmd, "timeout", "flag_timeout_desc")
	setFlagUsage(findPathCmd, "request-timeout", "flag_request_timeout_desc")

	LookCmd.Short = i18n.T("cmd_short_look")
	LookCmd.Long = i18n.T("cmd_long_look")

	googleCmd.Short = i18n.T("cmd_short_google")
	googleLinksCmd.Short = i18n.T("cmd_short_google_links")
	setFlagUsage(googleLinksCmd, "mode", "flag_google_mode_desc")

	apiCmd.Short = i18n.T("cmd_short_api")
	apiScanCmd.Short = i18n.T("cmd_short_api_scan")
	setFlagUsage(apiScanCmd, "spec", "flag_api_spec_desc")
	setFlagUsage(apiScanCmd, "host", "flag_api_host_desc")
	setFlagUsage(apiScanCmd, "show", "flag_api_show_desc")
	setFlagUsage(apiScanCmd, "exclude", "flag_api_exclude_desc")
	setFlagUsage(apiScanCmd, "request-timeout", "flag_request_timeout_desc")
	setFlagUsage(apiScanCmd, "jwt", "flag_api_jwt_desc")
	setFlagUsage(apiScanCmd, "api-key-header", "flag_api_key_header_desc")
	setFlagUsage(apiScanCmd, "api-key", "flag_api_key_desc")
	setFlagUsage(apiScanCmd, "spec-auth", "flag_api_spec_auth_desc")
	setFlagUsage(apiScanCmd, "show-errors", "flag_api_show_errors_desc")

	versionCmd.Short = i18n.T("cmd_short_version")
}

func setFlagUsage(cmd *cobra.Command, flagName string, msgKey string) {
	if cmd == nil {
		return
	}
	flag := cmd.Flags().Lookup(flagName)
	if flag == nil {
		flag = cmd.PersistentFlags().Lookup(flagName)
	}
	if flag != nil {
		flag.Usage = i18n.T(msgKey)
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
