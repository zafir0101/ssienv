package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
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
				fmt.Println("The command \"list\" is only available to institutional controllers")
				os.Exit(1)
			}

			ins := controller.(*domain.InstitutionController)
			if len(ins.PublishedDIDs) == 0 {
				fmt.Println("No stored published did on your controller")
				return
			}
			fmt.Printf("%-30s%-30s\n", "Label", "DID")
			for label, did := range ins.PublishedDIDs {
				fmt.Printf("%-30s%s\n", label, did)
			}
		},
	}
)
