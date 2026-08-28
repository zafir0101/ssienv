package create

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
)

var (
	connectionLabel string

	CreateCmd = &cobra.Command{
		Use:   "create",
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

			invOOB, err := controller.CreateConnection(connectionLabel)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("Send this invitation code to the other peer:\n %s", invOOB)

			if err := serializer.Serialize(controllerLabel, controller); err != nil {
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
