package create

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
)

var (
	offerLabel string
	connection string
	claims     string
	schema     string

	CreateCmd = &cobra.Command{
		Use:   "create",
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

			if err := ins.CreateCredentialOffer(offerLabel, json.RawMessage(claims), connection, schema); err != nil {
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
	CreateCmd.Flags().StringVarP(&offerLabel, "label", "l", "", "the label that will identifies the credential offer on your controller")
	CreateCmd.Flags().StringVarP(&connection, "connection", "", "", "The label (stored inside your controller) of a DIDComm connection that already "+
		"exists between the this issuer agent and the holder cloud or edeg agent (required)")
	CreateCmd.Flags().StringVarP(&claims, "claims", "", "", "The set of claims that will be included in the issued credential "+
		"(json raw) (required)")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "The URL pointing to the JSON schema that will be used for this offer "+
		"(should be 'http' or 'https') (required)")
}
