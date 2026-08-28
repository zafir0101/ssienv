package proof_request

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/proof-request/accept"
	"github.com/zafir0101/SSI-ENV/cmd/proof-request/create"
	"github.com/zafir0101/SSI-ENV/cmd/proof-request/list"
)

var (
	ProofRequestCmd = &cobra.Command{
		Use:   "proof-request",
		Short: "",
	}
)

func init() {
	ProofRequestCmd.AddCommand(create.CreateCmd)
	ProofRequestCmd.AddCommand(list.ListCmd)
	ProofRequestCmd.AddCommand(accept.AcceptCmd)
}
