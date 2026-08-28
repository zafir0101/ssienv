package whoami

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
)

var (
	WhoamiCmd = &cobra.Command{
		Use:   "whoami",
		Short: "Return the insitution/individual did",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithPureCommand(cmd, whoami); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func whoami(coData serializer.ControllerData) error {
	fmt.Println(coData.Controller.DID())
	return nil
}
