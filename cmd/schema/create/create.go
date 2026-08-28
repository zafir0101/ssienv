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
	schemaLabel string
	schema      string

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

			if err := ins.CreateSchema(schemaLabel, json.RawMessage(schema)); err != nil {
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
	CreateCmd.Flags().StringVarP(&schemaLabel, "label", "l", "", "the label that identifies the schema on your controller")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "the schema to be created (json raw)")
}
