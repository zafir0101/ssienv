package create

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
)

var (
	offerLabel string
	connection string
	claims     string
	schema     string

	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a credential offer and stored the record id on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, create); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&offerLabel, "label", "l", "", "the label that will identifies the credential offer on your controller")
	CreateCmd.Flags().StringVarP(&connection, "connection", "", "", "The label (stored inside your controller) of a DIDComm connection")
	CreateCmd.Flags().StringVarP(&claims, "claims", "", "", "The set of claims that will be included in the issued credential "+
		"(json raw) (required)")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "The URL pointing to the JSON schema that will be used for this offer "+
		"(should be 'http' or 'https') (required)")
}

func create(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"create\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	return ins.CreateCredentialOffer(offerLabel, json.RawMessage(claims), connection, schema)
}
