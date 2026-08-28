package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/connection"
	"github.com/zafir0101/ssienv/cmd/credential"
	"github.com/zafir0101/ssienv/cmd/did"
	init_controller "github.com/zafir0101/ssienv/cmd/init"
	"github.com/zafir0101/ssienv/cmd/proof-request"
	"github.com/zafir0101/ssienv/cmd/schema"
)

var (
	controllerLabel string

	rootHelpTemplate string = `ssienv is an identity environment that implements the basics concepts of SSI.
This application is a tool to generate and manage DIDs, 
Verifiable Credentials, and to handle Proof Requests.
Usage:
  ssienv [resource]
  ssienv [command]

Available Resources:
  did            Decentralized Identifier. Unique, user-controlled digital address
  connection     Peer-to-Peer DIDComm connection
  credential     Digital verifiable credential   
  proof-request  Cryptographic query sent by a verifier to a holder asking them to share specific attributes of their credential
  schema         Digital blueprint that defines the structure, attribute names, and data types for a verifiable credential     

Available Commands:
  init          Initialize and store the controller with the name and Url Agent API provided

Flags:
  -c, --controller string   Controller label (required for all subcommands of ssienv)
  -h, --help                Help for ssienv

Use "ssienv [command] --help" for more information about a command.
`

	rootCmd = &cobra.Command{
		Use: "ssienv",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.HasParent() {
				flag, err := cmd.Flags().GetString("controller")
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				if flag == "" {
					fmt.Println("Error: required flag \"controller\" not set\nType --help for more information")
					cmd.Help()
					os.Exit(1)
				}
			}
		},
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return nil
}

func init() {
	defaultHelpTemplate := rootCmd.HelpTemplate()
	rootCmd.SetHelpTemplate(rootHelpTemplate)

	rootCmd.AddCommand(init_controller.InitCmd)
	rootCmd.AddCommand(did.DIDCmd)
	rootCmd.AddCommand(credential.CredentialCmd)
	rootCmd.AddCommand(connection.ConnectionCmd)
	rootCmd.AddCommand(schema.SchemaCmd)
	rootCmd.AddCommand(proof_request.ProofRequestCmd)

	for _, cmd := range rootCmd.Commands() {
		if cmd == credential.CredentialCmd {
			continue
		}
		cmd.SetHelpTemplate(defaultHelpTemplate)
	}

	rootCmd.PersistentFlags().StringVarP(&controllerLabel, "controller", "c", "", "Controller label (required for all subcommands of ssienv)")
}
