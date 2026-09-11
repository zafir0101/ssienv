package create

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
)

var (
	nonDefault      bool
	proofReqLabel   string
	connectionLabel string
	schema          string

	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a proof request and store the presentation id on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, create); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&proofReqLabel, "label", "l", "", "The label that will identify the proof request on your controller (required)")
	CreateCmd.Flags().StringVarP(&connectionLabel, "connection", "", "", "The label that identifies the connection on your controller (required)")

	CreateCmd.Flags().BoolVar(&nonDefault, "non-default", false, "Use a schema other than the controller's default")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "The URL pointing to the JSON schema that will be used for this proof request (required if non-default is true)")

	CreateCmd.MarkFlagRequired("label")
	CreateCmd.MarkFlagRequired("connection")
	CreateCmd.MarkFlagsRequiredTogether("non-default", "schema")
}

func create(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"create\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	if nonDefault {
		return ins.CreateProofRequest(proofReqLabel, connectionLabel, schema)
	}

	return ins.CreateDefaultProofRequest(proofReqLabel, connectionLabel)
}
