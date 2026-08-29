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
	schemaLabel string
	schema      string

	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a schema and stored the id on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, create); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&schemaLabel, "label", "l", "", "the label that identifies the schema on your controller (required)")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "the schema to be created (json raw) (required)")

	CreateCmd.MarkFlagRequired("label")
	CreateCmd.MarkFlagRequired("schema")
}

func create(coData serializer.ControllerData) error {

	if !coData.IsInstitutional {
		return errors.New("The command \"publish\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	if err := ins.CreateSchema(schemaLabel, json.RawMessage(schema)); err != nil {
		return err
	}

	return nil
}
