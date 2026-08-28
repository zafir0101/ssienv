package create

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
)

var (
	connectionLabel string

	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a connection returning the invitation code",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, create); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&connectionLabel, "label", "l", "", "the label that identifies the connection on your controller (required)")

	CreateCmd.MarkFlagRequired("label")
}

func create(coData serializer.ControllerData) error {

	invOOB, err := coData.Controller.CreateConnection(connectionLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Send this invitation code to the other peer:\n %s", invOOB)
	return nil
}
