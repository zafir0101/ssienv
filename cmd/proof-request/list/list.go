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
	sent   bool
	accept bool

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
	ListCmd.Flags().BoolVarP(&accept, "accept", "a", false, "set the filter to list the offers accepted (only available to institutional controllers)")
}

func list(coData serializer.ControllerData) error {
	if !coData.IsInstitutional && sent {
		return errors.New("The flag \"sent\" is only available to institutional controllers")
	}

	if err := coData.Controller.RefreshProofRequests(); err != nil {
		return err
	}

	if sent && accept {
		ins := coData.Controller.(*domain.InstitutionController)

		proofReqs := ins.ProofRequestsSentAccepted
		if len(proofReqs) == 0 {
			fmt.Println("No stored proof request on your controller")
			return nil
		}

		fmt.Printf("%-30s%-30s\n", "Label", "Proof Request")
		for label, proofReq := range proofReqs {
			fmt.Printf("%-30s%s\n", label, proofReq)
		}

		return nil
	}

	if sent {
		ins := coData.Controller.(*domain.InstitutionController)

		proofReqs := ins.ProofRequestsSent
		if len(proofReqs) == 0 {
			fmt.Println("No stored proof request sent by your controller")
			return nil
		}

		fmt.Printf("%-30s%-30s\n", "Label", "Proof Request")
		for label, proofReq := range proofReqs {
			fmt.Printf("%-30s%s\n", label, proofReq)
		}

		return nil
	}

	if accept {
		proofReqs := coData.Controller.ProofRequestsAcceptedHashMap()
		if len(proofReqs) == 0 {
			fmt.Println("No stored proof request accepted by your controller")
			return nil
		}

		fmt.Printf("%-30s%-30s\n", "Label", "Proof Request")
		for label, proofReq := range proofReqs {
			fmt.Printf("%-30s%s\n", label, proofReq)
		}

		return nil
	}

	proofReqs := coData.Controller.ProofRequestsReceivedHashMap()
	if len(proofReqs) == 0 {
		fmt.Println("No stored proof request received on your controller")
		return nil
	}

	fmt.Printf("Proof Request\n")
	for _, proofReq := range proofReqs {
		fmt.Printf("%s\n", proofReq)
	}

	return nil
}
