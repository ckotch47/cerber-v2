package command

import (
	"fmt"

	"cerber/internal/i18n"
	"cerber/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: i18n.T("cmd_short_version"),
	RunE:  getVersion,
}

func getVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(i18n.T("msg_version", version.String))
	return nil
}
