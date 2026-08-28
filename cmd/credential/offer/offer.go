package offer

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/credential/offer/accept"
	"github.com/zafir0101/SSI-ENV/cmd/credential/offer/create"
	"github.com/zafir0101/SSI-ENV/cmd/credential/offer/list"
)

var (
	OfferCmd = &cobra.Command{
		Use:   "offer",
		Short: "",
	}
)

func init() {
	OfferCmd.AddCommand(create.CreateCmd)
	OfferCmd.AddCommand(list.ListCmd)
	OfferCmd.AddCommand(accept.AcceptCmd)
}
