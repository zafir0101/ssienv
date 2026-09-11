package create

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
	nonDefault bool
	offerLabel string
	connection string
	claims     string
	schema     string

	// Flags do esquema brEduPerson (default)
	cpf             string
	affiliationID   string
	affiliationType string
	entranceDate    string
	exitDate        string

	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a credential offer and store the record id on your controller",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serializer.WithMutateCommand(cmd, create); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	CreateCmd.Flags().StringVarP(&offerLabel, "label", "l", "", "The label that will identifies the credential offer on your controller")
	CreateCmd.Flags().StringVarP(&connection, "connection", "", "", "The label (stored inside your controller) of a DIDComm connection")

	// Flags para o caso non-default
	CreateCmd.Flags().StringVarP(&claims, "claims", "", "", "The set of claims that will be included in the issued credential (json raw) (required if non-default is true)")
	CreateCmd.Flags().StringVarP(&schema, "schema", "s", "", "The URL pointing to the JSON schema that will be used for this offer (should be 'http' or 'https') (required if non-default is true)")
	CreateCmd.Flags().BoolVarP(&nonDefault, "non-default", "", false, "Indicates whether the credential will follow the pattern defined by the controller")

	// Flags para o caso default (esquema brEduPerson)
	CreateCmd.Flags().StringVarP(&cpf, "cpf", "", "", "The subject's CPF (brPersonCPF) (required if using default schema)")
	CreateCmd.Flags().StringVarP(&affiliationID, "affiliation-id", "", "", "Identifier of this affiliation instance (brEduAffiliation) (required if using default schema)")
	CreateCmd.Flags().StringVarP(&affiliationType, "affiliation-type", "r", "", "The type of affiliation with the institution: student, employee, staff, faculty, etc (brEduAffiliationType) (required if using default schema)")
	CreateCmd.Flags().StringVarP(&entranceDate, "entrance-date", "e", "", "Date the affiliation began, format YYYYMMDD (brEntranceDate)")
	CreateCmd.Flags().StringVarP(&exitDate, "exit-date", "x", "", "Date the affiliation ended, format YYYYMMDD (brExitDate) (optional)")

	CreateCmd.MarkFlagsRequiredTogether("non-default", "claims", "schema")
	CreateCmd.MarkFlagsMutuallyExclusive("non-default", "cpf")
	CreateCmd.MarkFlagsMutuallyExclusive("non-default", "affiliation-id")
	CreateCmd.MarkFlagsMutuallyExclusive("non-default", "affiliation-type")
	CreateCmd.MarkFlagsRequiredTogether("cpf", "affiliation-id", "affiliation-type")
}

func create(coData serializer.ControllerData) error {
	if !coData.IsInstitutional {
		return errors.New("the command \"create\" is only available to institutional controllers")
	}

	var finalClaims json.RawMessage
	var finalSchema string
	ins := coData.Controller.(*domain.InstitutionController)

	if nonDefault {
		if !json.Valid([]byte(claims)) {
			return errors.New("the provided claims string is not a valid JSON")
		}

		finalClaims = json.RawMessage(claims)
		finalSchema = schema

		return ins.CreateCredentialOffer(offerLabel, finalClaims, connection, finalSchema)
	}

	defaultClaims := map[string]string{
		"brPersonCPF":          cpf,
		"brEduAffiliation":     affiliationID,
		"brEduAffiliationType": affiliationType,
	}

	if entranceDate != "" {
		defaultClaims["brEntranceDate"] = entranceDate
	}
	if exitDate != "" {
		defaultClaims["brExitDate"] = exitDate
	}

	claimsBytes, err := json.Marshal(defaultClaims)
	if err != nil {
		return fmt.Errorf("failed to encode default claims: %v", err)
	}

	finalClaims = json.RawMessage(claimsBytes)

	return ins.CreateDefaultCredentialOffer(offerLabel, finalClaims, connection)
}
