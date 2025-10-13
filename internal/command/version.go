package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Показать версию приложения",
	Run: getVersion,
}

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Println("Версия: v0.0.1a")
}