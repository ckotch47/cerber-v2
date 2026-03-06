package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"cerber/internal/recon"
	"cerber/internal/style"
)

var googleMode string

var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "Google dork helpers",
}

var googleLinksCmd = &cobra.Command{
	Use:   "links <domain>",
	Short: "Сгенерировать dork-ссылки для домена",
	Args:  cobra.ExactArgs(1),
	RunE:  runGoogleLinks,
}

func init() {
	googleLinksCmd.Flags().StringVar(
		&googleMode,
		"mode",
		"all",
		"Режимы через запятую (например: 1,5,12) или all",
	)
	googleCmd.AddCommand(googleLinksCmd)
}

func runGoogleLinks(_ *cobra.Command, args []string) error {
	domain, err := cleanDomain(args[0])
	if err != nil {
		return err
	}

	service := recon.NewGoogleLinksService()
	links, err := service.GenerateLinks(domain, googleMode)
	if err != nil {
		return err
	}

	for _, link := range links {
		fmt.Println(style.SuccessStyle.Render(link))
	}
	return nil
}
