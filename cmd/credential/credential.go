package credential

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/credential/list"
	"github.com/zafir0101/ssienv/cmd/credential/offer"
)

var (
	helpTemplate string = `Digital verifiable credential

Usage:
  ssienv credential [command]

Available Resources:
  offer       Digital message sent by an issuer to a holder that starts the process of issuing a verifiable credential

Available Commands:
  list        List the stored credentials on your controller

Flags:
  -h, --help   help for credential

Global Flags:
  -c, --controller string   Controller label (required for all subcommands of ssienv)

Use "ssienv credential [command] --help" for more information about a command.
`
	CredentialCmd = &cobra.Command{
		Use:   "credential",
		Short: "Digital verifiable credential",
	}
)

func init() {
	defaultHelpTemplate := CredentialCmd.HelpTemplate()

	offer.OfferCmd.SetHelpTemplate(defaultHelpTemplate)
	CredentialCmd.SetHelpTemplate(helpTemplate)

	CredentialCmd.AddCommand(offer.OfferCmd)
	CredentialCmd.AddCommand(list.ListCmd)

}
