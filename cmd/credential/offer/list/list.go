package list

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
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
	if !coData.IsInstitutional && sent {
		return errors.New("The flag \"sent\" is only available to institutional controllers")
	}

	if sent {
		ins := coData.Controller.(*domain.InstitutionController)

		offers := ins.CredentialOffersSent
		if len(offers) == 0 {
			fmt.Println("No stored offer on your controller")
			return nil
		}

		fmt.Printf("%-30s%-30s\n", "Label", "Offers")
		for label, offer := range offers {
			fmt.Printf("%-30s%s\n", label, offer)
		}
		return nil
	}

	if err := coData.Controller.RefreshOffersReceived(); err != nil {
		return err
	}

	offers := coData.Controller.CredentialOffersReceivedSlice()
	if len(offers) == 0 {
		fmt.Println("No stored offer on your controller")
		return nil
	}

	fmt.Printf("Offers\n")
	for _, offer := range offers {
		fmt.Printf("%s\n", offer)
	}

	return nil
}
