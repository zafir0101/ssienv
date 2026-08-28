package create

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
)

var (
	proofReqLabel   string
	connectionLabel string
	schema          string

	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, create); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&proofReqLabel, "label", "l", "", "the label that will identifies the proof request on your controller (required)")
	CreateCmd.Flags().StringVarP(&connectionLabel, "connection", "", "", "the label that identifies the connection on your controller (required)")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "The URL pointing to the JSON schema that will be used for this offer (required)")

	CreateCmd.MarkFlagRequired("label")
	CreateCmd.MarkFlagRequired("connection")
}

func create(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"create\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	if err := ins.CreateProofRequest(proofReqLabel, connectionLabel, schema); err != nil {
		return err
	}

	return nil
}
