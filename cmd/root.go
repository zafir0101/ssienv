package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/connection"
	"github.com/zafir0101/SSI-ENV/cmd/credential"
	"github.com/zafir0101/SSI-ENV/cmd/did"
	proof_request "github.com/zafir0101/SSI-ENV/cmd/proof-request"
	"github.com/zafir0101/SSI-ENV/cmd/schema"
)

var (
	controllerLabel string

	rootCmd = &cobra.Command{
		Use:   "ssienv",
		Short: "ssienv is an identity environment that implements the basics concepts of SSI",
		Long: `ssienv is an identity environment that implements the basics concepts of SSI.
This application is a tool to gerenate and manage DIDs,
Verifiable Credentials and handle Proofs Requests.`,
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
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(did.DIDCmd)
	rootCmd.AddCommand(credential.CredentialCmd)
	rootCmd.AddCommand(connection.ConnectionCmd)
	rootCmd.AddCommand(schema.SchemaCmd)
	rootCmd.AddCommand(proof_request.ProofRequestCmd)

	rootCmd.PersistentFlags().StringVarP(&controllerLabel, "controller", "c", "", "controller label (required for all subcommands of ssienv)")
}
