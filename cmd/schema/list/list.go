package list

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List the stored schemas on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithPureCommand(cmd, list); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func list(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"list\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)
	if len(ins.Schemas) == 0 {
		fmt.Println("No stored schemas on your controller")
		return nil
	}

	fmt.Printf("%-30s%-30s\n", "Label", "SchemaID")
	for label, schema := range ins.Schemas {
		fmt.Printf("%-30s%s\n", label, schema)
	}

	return nil
}
