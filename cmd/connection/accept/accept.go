package accept

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
)

var (
	connectionLabel string
	invOOB          string

	AcceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			controllerLabel, err := cmd.Flags().GetString("controller")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			controller, _, err := serializer.Deserialize(controllerLabel)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if err := controller.AcceptConnection(connectionLabel, invOOB); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if err := serializer.Serialize(controllerLabel, controller); err != nil {
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
