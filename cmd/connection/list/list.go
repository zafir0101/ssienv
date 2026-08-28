package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List the stored connections on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithPureCommand(cmd, list); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func list(coData serializer.ControllerData) error {
	connections := coData.Controller.ConnectionsHashMap()
	if len(connections) == 0 {
		fmt.Println("No stored connections on your controller")
		return nil
	}

	fmt.Printf("%-30s%-30s\n", "Label", "Connections")
	for label, did := range connections {
		fmt.Printf("%-30s%s\n", label, did)
	}

	return nil
}
