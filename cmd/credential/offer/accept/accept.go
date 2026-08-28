package accept

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
)

var (
	credentialLabel string
	recordID        string

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
	AcceptCmd.Flags().StringVarP(&credentialLabel, "label", "l", "", "the label that will identifies the credential on your controller (required)")
	AcceptCmd.Flags().StringVarP(&recordID, "record-id", "r", "", "the unique identifier of the issue credential record (required)")

	AcceptCmd.MarkFlagRequired("label")
	AcceptCmd.MarkFlagRequired("record-id")
}

func accept(coData serializer.ControllerData) error {
	return coData.Controller.AcceptCredentialOffer(credentialLabel, recordID)
}
