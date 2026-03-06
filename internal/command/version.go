package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Показать версию приложения",
	RunE:  getVersion,
}

func getVersion(cmd *cobra.Command, args []string) error {
	fmt.Println("Версия: v0.0.1")
	return nil
}
