package resolve

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
)

var (
	did string

	ResolveCmd = &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a did (short-form) and return a did document",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithPureCommand(cmd, resolve); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	ResolveCmd.Flags().StringVarP(&did, "did", "d", "", "did short form to be resolved (required)")

	ResolveCmd.MarkFlagRequired("did")
}

func resolve(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"resolve\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	didDocument, err := ins.ResolveDID(did)
	if err != nil {
		return err
	}

	json, err := json.MarshalIndent(didDocument, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}
