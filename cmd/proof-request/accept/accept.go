package accept

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
)

var (
	proofReqLabel   string
	credentialLabel string
	presID          string

	AcceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, accept); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	AcceptCmd.Flags().StringVarP(&proofReqLabel, "label", "l", "", "the label that will identifies the proof request on your controller (required)")
	AcceptCmd.Flags().StringVarP(&credentialLabel, "credential", "", "", "the label that identifies the credential on your controller (required)")
	AcceptCmd.Flags().StringVarP(&presID, "presentation-id", "p", "", "the unique identifier of the proof request (required)")

	AcceptCmd.MarkFlagRequired("label")
	AcceptCmd.MarkFlagRequired("credential")
	AcceptCmd.MarkFlagRequired("presentation-id")
}

func accept(coData serializer.ControllerData) error {
	return coData.Controller.AcceptProofRequest(proofReqLabel, credentialLabel, presID)
}
