package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
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

			controller, _, err := serializer.Deserialize(controllerLabel)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			connections := controller.ConnectionsHashMap()
			if len(connections) == 0 {
				fmt.Println("No stored connection on your controller")
				return
			}

			fmt.Printf("%-30s%-30s\n", "Label", "Connections")
			for label, did := range connections {
				fmt.Printf("%-30s%s\n", label, did)
			}
		},
	}
)
