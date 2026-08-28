package publish

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
)

var (
	didLabel string
	did      string

	PublishCmd = &cobra.Command{
		Use:   "publish",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			controllerLabel, err := cmd.Flags().GetString("controller")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			controller, isInstitutional, err := serializer.Deserialize(controllerLabel)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if !isInstitutional {
				fmt.Println("The command \"publish\" is only available to institutional controllers")
				os.Exit(1)
			}

			ins := controller.(*domain.InstitutionController)

			if err := ins.PublishDID(didLabel, did); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if err := serializer.Serialize(controllerLabel, ins); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	PublishCmd.Flags().StringVarP(&didLabel, "label", "l", "", "the label that identifies the did on your controller")
	PublishCmd.Flags().StringVarP(&did, "did", "d", "", "did long form (required)")

	PublishCmd.MarkFlagRequired("did")
}
