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
		Short: "List the stored credentials on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithPureCommand(cmd, list); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {}

func list(coData serializer.ControllerData) error {
	creds := coData.Controller.CredentialsHashMap()
	if len(creds) == 0 {
		fmt.Println("No stored credentials on your controller")
		return nil
	}

	fmt.Printf("%-30s%-30s\n", "Label", "Credential")
	for label, cred := range creds {
		fmt.Printf("%-30s%s\n", label, cred)
	}

	return nil
}
