package offer

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/credential/offer/accept"
	"github.com/zafir0101/ssienv/cmd/credential/offer/create"
	"github.com/zafir0101/ssienv/cmd/credential/offer/list"
)

var (
	OfferCmd = &cobra.Command{
		Use:   "offer",
		Short: "Digital message sent by an issuer to a holder that starts the process of issuing a verifiable credential",
	}
)

func init() {
	OfferCmd.AddCommand(create.CreateCmd)
	OfferCmd.AddCommand(list.ListCmd)
	OfferCmd.AddCommand(accept.AcceptCmd)
}
