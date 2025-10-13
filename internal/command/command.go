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

func cleanDomain(domain string) string {
	// Убираем "http://", "https://", "www."
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimSuffix(domain, "/")
	
	if res := strings.HasPrefix(domain, "www."); res {
		domain = domain[4:]
	}
	return domain
}
