package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

var (
	agentURL      string
	institutional bool

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize and store the controller with the name and Url Agent API provided",
		Run: func(cmd *cobra.Command, args []string) {
			controllerLabel, err := cmd.Flags().GetString("controller")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			parsedURL, err := url.Parse(agentURL)
			if err != nil {
				fmt.Printf("An error occured while parsing the agent url: %s\n", err.Error())
				os.Exit(1)
			}

			if institutional {
				api := ssi.NewCloudAgentAPI(parsedURL)
				controller := domain.NewInstitutionController(api)
				if err := serializer.Serialize(controllerLabel, controller); err != nil {
					fmt.Printf("An error occured while serializing the controller: %s\n", err.Error())
					os.Exit(1)
				}
			} else {
				api := ssi.NewEdgeAgentAPI(parsedURL)
				controller := domain.NewIndividualController(api)
				if err := serializer.Serialize(controllerLabel, controller); err != nil {
					fmt.Printf("An error occured while serializing the controller: %s\n", err.Error())
					os.Exit(1)
				}
			}
		},
	}
)

func init() {
	initCmd.Flags().StringVarP(&controllerLabel, "controller", "c", "", "controller label (required)")
	initCmd.Flags().StringVarP(&agentURL, "url", "u", "", "agent url (required)")
	initCmd.Flags().BoolVarP(&institutional, "institutional", "i", false, "whether the controller is to an institution")

	initCmd.MarkFlagRequired("url")
}
