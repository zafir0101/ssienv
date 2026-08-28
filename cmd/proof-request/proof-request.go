package proof_request

import (
	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/proof-request/accept"
	"github.com/zafir0101/ssienv/cmd/proof-request/create"
	"github.com/zafir0101/ssienv/cmd/proof-request/list"
)

var (
	ProofRequestCmd = &cobra.Command{
		Use:   "proof-request",
		Short: "Cryptographic query sent by a verifier to a holder asking them to share specific attributes of their credential",
	}
)

func init() {
	ProofRequestCmd.AddCommand(create.CreateCmd)
	ProofRequestCmd.AddCommand(list.ListCmd)
	ProofRequestCmd.AddCommand(accept.AcceptCmd)
}
