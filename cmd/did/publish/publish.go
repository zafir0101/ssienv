package publish

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
)

var (
	didLabel string
	did      string

	PublishCmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish long-form did",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, publish); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	PublishCmd.Flags().StringVarP(&didLabel, "label", "l", "", "the label that identifies the did on your controller")
	PublishCmd.Flags().StringVarP(&did, "did", "d", "", "long-form did (required)")

	PublishCmd.MarkFlagRequired("did")
}

func publish(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"publish\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	if err := ins.PublishDID(didLabel, did); err != nil {
		return err
	}

	return nil
}
