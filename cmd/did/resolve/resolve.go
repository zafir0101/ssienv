package resolve

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
)

var (
	did string

	ResolveCmd = &cobra.Command{
		Use:   "resolve",
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
				fmt.Println("The command \"resolve\" is only available to institutional controllers")
				os.Exit(1)
			}

			ins := controller.(*domain.InstitutionController)

			didDocument, err := ins.ResolveDID(did)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			json, err := json.MarshalIndent(didDocument, "", " ")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Println(string(json))

		},
	}
)

func init() {
	ResolveCmd.Flags().StringVarP(&did, "did", "d", "", "did short form to be resolved (required)")

	ResolveCmd.MarkFlagRequired("did")
}
