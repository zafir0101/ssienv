package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
)

var (
	sent bool

	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithPureCommand(cmd, list); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	ListCmd.Flags().BoolVarP(&sent, "sent", "s", false, "set the filter to list the offers sent (only available to institutional controllers)")
}

func list(coData serializer.ControllerData) error {
	creds := coData.Controller.CredentialsHashMap()
	if len(creds) == 0 {
		fmt.Println("No stored credential on your controller")
		return nil
	}

	fmt.Printf("%-30s%-30s\n", "Label", "Credential")
	for label, cred := range creds {
		fmt.Printf("%-30s%s\n", label, cred)
	}

	return nil
}
