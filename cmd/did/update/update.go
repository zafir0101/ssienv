package update

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zafir0101/SSI-ENV/cmd/serializer"
	"github.com/zafir0101/SSI-ENV/internal/domain"
)

var (
	keyID      string
	keyPurpose int
	remove     bool

	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			controllerLabel, err := cmd.Flags().GetString("controller")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			controller, isInstitutional, err := serializer.Deserialize(controllerLabel)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if !isInstitutional {
				fmt.Println("The command \"update\" is only available to institutional controllers")
				os.Exit(1)
			}

			ins := controller.(*domain.InstitutionController)

			if remove {
				if err := ins.RemoveDIDKey(keyID, domain.KeyPurpose(keyPurpose)); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			} else {
				if err := ins.AddKeyToDID(domain.KeyPurpose(keyPurpose)); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}

			if err := serializer.Serialize(controllerLabel, ins); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	UpdateCmd.Flags().StringVarP(&keyID, "id", "i", "", "id key (only for add operation)")
	UpdateCmd.Flags().IntVarP(&keyPurpose, "purpose", "p", 0, "key purpose")
	UpdateCmd.Flags().BoolVarP(&remove, "remove", "r", false, "select the remove key operation.")

	UpdateCmd.MarkFlagsRequiredTogether("remove", "id")
}
