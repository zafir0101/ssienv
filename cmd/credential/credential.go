package credential

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/credential/list"
	"github.com/zafir0101/SSI-ENV/cmd/credential/offer"
)

var (
	CredentialCmd = &cobra.Command{
		Use:   "credential",
		Short: "",
	}
)

func init() {
	CredentialCmd.AddCommand(offer.OfferCmd)
	CredentialCmd.AddCommand(list.ListCmd)

}
