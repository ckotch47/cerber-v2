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
	Short: "Краткое описание",
	Long:  `Полное описание моего приложения`,
}


func init(){
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(versionCmd)
	
	rootCmd.AddCommand(findCmd)
	findCmd.AddCommand(findAdminCmd)

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
