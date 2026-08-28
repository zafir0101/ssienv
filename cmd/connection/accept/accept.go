package accept

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
)

var (
	connectionLabel string
	invOOB          string

	AcceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accept a connection invitation received by your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, accept); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	AcceptCmd.Flags().StringVarP(&connectionLabel, "label", "l", "", "the label that identifies the connection on your controller (required)")
	AcceptCmd.Flags().StringVarP(&invOOB, "invitation", "i", "", "the invitation code (required)")

	AcceptCmd.MarkFlagRequired("label")
	AcceptCmd.MarkFlagRequired("invitation")
}

func accept(coData serializer.ControllerData) error {
	return coData.Controller.AcceptConnection(connectionLabel, invOOB)
}
