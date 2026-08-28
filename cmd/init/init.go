package init_controller

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
	"github.com/zafir0101/ssienv/internal/ssi"
)

var (
	agentURL      string
	institutional bool

	InitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize and store the controller with the provided name and Agent API URL",
		Run: func(cmd *cobra.Command, args []string) {
			controllerLabel, err := cmd.Flags().GetString("controller")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			parsedURL, err := url.Parse(agentURL)
			if err != nil {
				fmt.Printf("An error occurred while parsing the agent url: %s\n", err.Error())
				os.Exit(1)
			}

			if institutional {
				api := ssi.NewCloudAgentAPI(parsedURL)
				controller := domain.NewInstitutionController(api)
				if err := serializer.Serialize(controllerLabel, controller); err != nil {
					fmt.Printf("An error occurred while serializing the controller: %s\n", err.Error())
					os.Exit(1)
				}
			} else {
				api := ssi.NewEdgeAgentAPI(parsedURL)
				controller := domain.NewIndividualController(api)
				if err := serializer.Serialize(controllerLabel, controller); err != nil {
					fmt.Printf("An error occurred while serializing the controller: %s\n", err.Error())
					os.Exit(1)
				}
			}
		},
	}
)

func init() {
	InitCmd.Flags().StringVarP(&agentURL, "url", "u", "", "Agent url (required)")
	InitCmd.Flags().BoolVarP(&institutional, "institutional", "i", false, "Whether the controller is to an institution")

	InitCmd.MarkFlagRequired("url")
}
