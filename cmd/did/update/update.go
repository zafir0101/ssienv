package update

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/ssienv/cmd/serializer"
	"github.com/zafir0101/ssienv/internal/domain"
)

var (
	keyID      string
	keyPurpose int
	remove     bool

	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the keys of institutional/individual did",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, update); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVarP(&keyID, "id", "i", "", "key ID to remove (only for the remove operation)")
	UpdateCmd.Flags().IntVarP(&keyPurpose, "purpose", "p", 0, "key purpose")
	UpdateCmd.Flags().BoolVarP(&remove, "remove", "r", false, "select the remove key operation")

	UpdateCmd.MarkFlagsRequiredTogether("remove", "id")
}

func update(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("The command \"update\" is only available to institutional controllers")
	}

	ins := coData.Controller.(*domain.InstitutionController)

	if remove {
		return ins.RemoveDIDKey(keyID, domain.KeyPurpose(keyPurpose))
	}

	return ins.AddKeyToDID(domain.KeyPurpose(keyPurpose))
}
