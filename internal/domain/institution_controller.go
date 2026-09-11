package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/zafir0101/ssienv/internal/ssi"
)

const defaultSchemaRaw = `{
    "$id": "https://example.com/breduperson-1.0.0",
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
      "brPersonCPF": { "type": "string" },
      "brEduAffiliation": { "type": "string" },
      "brEduAffiliationType": {
        "type": "string",
        "enum": ["student", "employee", "staff", "faculty"]
      },
      "brEntranceDate": { "type": "string" },
      "brExitDate": { "type": "string" }
    },
    "required": ["brPersonCPF", "brEduAffiliation", "brEduAffiliationType"],
    "additionalProperties": false
  }`

type InstitutionController struct {
	CloudAgentAPI *ssi.CloudAgentAPI

	InstitutionDIDPrism ssi.DIDPrism
	PublishedDIDs       map[string]ssi.DIDPrism

	Connections map[string]ssi.ConnectionID

	Schemas         map[string]ssi.SchemaID
	DefaultSchemaID string

	Credentials              map[string]ssi.RecordID
	CredentialOffersReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a oferta
	CredentialOffersSent     map[string]ssi.RecordID

	ProofRequestsAccepted     map[string]ssi.PresentationID
	ProofRequestsSentAccepted map[string]ssi.RecordID
	ProofRequestsReceived     []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a requisicao
	ProofRequestsSent         map[string]ssi.PresentationID

	Num_keys int
}

func NewInstitutionController(cloudAgentAPI *ssi.CloudAgentAPI) (*InstitutionController, error) {
	controller := &InstitutionController{
		CloudAgentAPI:             cloudAgentAPI,
		PublishedDIDs:             make(map[string]ssi.DIDPrism),
		Connections:               make(map[string]ssi.ConnectionID),
		Schemas:                   make(map[string]ssi.SchemaID),
		Credentials:               make(map[string]ssi.RecordID),
		CredentialOffersSent:      make(map[string]ssi.RecordID),
		ProofRequestsAccepted:     make(map[string]ssi.PresentationID),
		ProofRequestsSent:         make(map[string]ssi.PresentationID),
		ProofRequestsSentAccepted: make(map[string]ssi.RecordID),
	}

	if err := controller.createDID(); err != nil {
		return nil, err
	}

	if err := controller.CreateSchema("defaultSchema", json.RawMessage(defaultSchemaRaw)); err != nil {
		return nil, err
	}
	controller.DefaultSchemaID = controller.Schemas["defaultSchema"]

	return controller, nil

}

func (co *InstitutionController) RefreshOffersReceived() error {
	recordIDs, credentialStatus, err := co.CloudAgentAPI.ListCredentialOffers()
	if err != nil {
		return err
	}

	clear(co.CredentialOffersReceived)
	co.CredentialOffersReceived = co.CredentialOffersReceived[:0]

	for i, credStatus := range credentialStatus {
		if credStatus == "OfferReceived" {
			co.CredentialOffersReceived = append(co.CredentialOffersReceived, recordIDs[i])
		}
	}

	return nil
}

func (co *InstitutionController) RefreshProofRequests() error {
	presentationIDs, proofReqStatus, err := co.CloudAgentAPI.ListProofRequestsData()
	if err != nil {
		return err
	}

	clear(co.ProofRequestsReceived)
	co.ProofRequestsReceived = co.ProofRequestsReceived[:0]

	for i, proofReqStatus := range proofReqStatus {
		if proofReqStatus == "RequestReceived" {
			co.ProofRequestsReceived = append(co.ProofRequestsReceived, presentationIDs[i])
			continue
		}

		if proofReqStatus == "PresentationVerified" {
			for label, preID := range co.ProofRequestsSent {
				if presentationIDs[i] == preID {
					co.ProofRequestsSentAccepted[label] = preID
					delete(co.ProofRequestsSent, label)
					break
				}
			}
		}
	}

	return nil
}

func (co *InstitutionController) createDID() error {
	pksID := []string{"key1-authentication", "key2-assertionMethod"}
	pksPurpose := []KeyPurpose{Authentication, AssertionMethod}

	payload, err := newDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return err
	}

	did, err := co.CloudAgentAPI.CreateDID(payload)
	if err != nil {
		return err
	}

	co.InstitutionDIDPrism = did
	co.Num_keys = 2

	return nil
}

func (co *InstitutionController) PublishDID(label string, didLongForm ssi.LongFormDIDPrism) error {
	did, err := co.CloudAgentAPI.PublishDID(didLongForm)
	if err != nil {
		return err
	}

	co.PublishedDIDs[label] = did
	return nil
}

func (co *InstitutionController) ResolveDID(did ssi.DIDPrism) (ssi.DIDPrismDocument, error) {
	didDoc, err := co.CloudAgentAPI.ResolveDID(did)
	if err != nil {
		return ssi.DIDPrismDocument{}, err
	}

	return didDoc, err
}

func (co *InstitutionController) AddKeyToDID(pkPurpose KeyPurpose) error {
	if !pkPurpose.isValid() {
		return errors.New("invalid key purpose")
	}

	pkID := "key" + strconv.Itoa(co.Num_keys+1) + "-" + pkPurpose.string()
	payload, err := newDIDUpdatePayload(addKey, pkID, pkPurpose)
	if err != nil {
		return err
	}

	if err := co.CloudAgentAPI.UpdateDID(payload, co.InstitutionDIDPrism); err != nil {
		return err
	}

	co.Num_keys++
	fmt.Println(co.Num_keys)
	return nil
}

func (co *InstitutionController) RemoveDIDKey(pkID string, pkPurpose KeyPurpose) error {
	if !pkPurpose.isValid() {
		return errors.New("invalid key purpose")
	}

	payload, err := newDIDUpdatePayload(removeKey, pkID, pkPurpose)
	if err != nil {
		return err
	}

	if err := co.CloudAgentAPI.UpdateDID(payload, co.InstitutionDIDPrism); err != nil {
		return err
	}

	co.Num_keys--

	return nil
}

func (co *InstitutionController) deactivateDID() error {
	return co.CloudAgentAPI.DeactivateDID(co.InstitutionDIDPrism)
}

func (co *InstitutionController) CreateConnection(label string) (ssi.InvitationOOB, error) {
	payload := newConnectionCreationPayload(label)

	connID, invOOB, err := co.CloudAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", err
	}

	co.Connections[label] = connID
	return invOOB, nil
}

func (co *InstitutionController) AcceptConnection(label string, invOOB ssi.InvitationOOB) error {
	payload := newConnectionAcceptPayload(invOOB)

	connID, err := co.CloudAgentAPI.AcceptConnection(payload)
	if err != nil {
		return err
	}

	co.Connections[label] = connID
	return nil
}

func (co *InstitutionController) deactivateConnection(label string) error {
	connID := co.Connections[label]
	if connID == "" {
		return errors.New("No Connections with label " + label)
	}

	if err := co.CloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	delete(co.Connections, label)

	return nil
}

func (co *InstitutionController) CreateSchema(schemaLabel string, schema json.RawMessage) error {
	payload := newSchemaCreationPayload(schemaLabel, co.InstitutionDIDPrism, schema)

	schemaID, err := co.CloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return err
	}

	co.Schemas[schemaLabel] = schemaID

	return nil
}

func (co *InstitutionController) CreateCredentialOffer(offerLabel string, claims json.RawMessage,
	connLabel string, schemaID ssi.SchemaID) error {
	connID := co.Connections[connLabel]
	if connID == "" {
		return errors.New("No Connections with label " + connLabel)
	}

	payload := newCredentialOfferPayload(claims, co.InstitutionDIDPrism, connID, schemaID)

	recordID, err := co.CloudAgentAPI.CreateCredentialOffer(payload)
	if err != nil {
		return err
	}

	co.CredentialOffersSent[offerLabel] = recordID

	return nil
}

func (co *InstitutionController) CreateDefaultCredentialOffer(offerLabel string, claims json.RawMessage,
	connLabel string) error {
	connID := co.Connections[connLabel]
	if connID == "" {
		return errors.New("No Connections with label " + connLabel)
	}

	payload := newCredentialOfferPayload(claims, co.InstitutionDIDPrism, connID, co.DefaultSchemaID)

	recordID, err := co.CloudAgentAPI.CreateCredentialOffer(payload)
	if err != nil {
		return err
	}

	co.CredentialOffersSent[offerLabel] = recordID

	return nil
}

func (co *InstitutionController) AcceptCredentialOffer(credentialLabel string, recID ssi.RecordID) error {
	payload := newOfferAcceptancePayload(co.InstitutionDIDPrism)

	for i, offer := range co.CredentialOffersReceived {
		if offer == recID {
			co.CredentialOffersReceived = append(co.CredentialOffersReceived[:i], co.CredentialOffersReceived[i+1:]...)
			break
		}

		if i == len(co.CredentialOffersReceived)-1 {
			return errors.New("No credential offer with this record id on your controller")
		}
	}

	if err := co.CloudAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	co.Credentials[credentialLabel] = recID

	return nil
}

// O Identus possui uma limitação: O holder sempre recebe uma requisição aberta, não sabendo exatamente
// o que deve apresentar. Isso porque o campo proof é settado para empty na implementação do identus. Não
// sei exatamente se o dado do campo proof é usado internamente para validar a apresentação do holder.
func (co *InstitutionController) CreateProofRequest(proofReqLabel string, connLabel string,
	schemaID ssi.SchemaID) error {
	connID := co.Connections[connLabel]
	if connID == "" {
		return errors.New("No Connections with label " + connLabel)
	}

	payload := newProofRequestPayload(proofReqLabel, connID, schemaID, co.InstitutionDIDPrism)

	presentationID, err := co.CloudAgentAPI.CreateProofRequest(payload)
	if err != nil {
		return err
	}

	co.ProofRequestsSent[proofReqLabel] = presentationID

	return nil
}

func (co *InstitutionController) CreateDefaultProofRequest(proofReqLabel string, connLabel string) error {
	connID := co.Connections[connLabel]
	if connID == "" {
		return errors.New("No Connections with label " + connLabel)
	}

	payload := newProofRequestPayload(proofReqLabel, connID, co.DefaultSchemaID, co.InstitutionDIDPrism)

	presentationID, err := co.CloudAgentAPI.CreateProofRequest(payload)
	if err != nil {
		return err
	}

	co.ProofRequestsSent[proofReqLabel] = presentationID

	return nil
}

func (co *InstitutionController) AcceptProofRequest(proofReqLabel string, credentialLabel string, presID ssi.PresentationID) error {
	recID := co.Credentials[credentialLabel]
	if recID == "" {
		return errors.New("No credentials with label " + credentialLabel)
	}

	for i, request := range co.ProofRequestsReceived {
		if request == presID {
			co.ProofRequestsReceived = append(co.ProofRequestsReceived[:i], co.ProofRequestsReceived[i+1:]...)
			break
		}

		if i == len(co.ProofRequestsReceived)-1 {
			return errors.New("No proof request with this record id on your controller")
		}
	}

	payload := newProofRequestAcceptancePayload(recID)

	if err := co.CloudAgentAPI.AcceptProofRequest(payload, presID); err != nil {
		return err
	}

	co.ProofRequestsAccepted[proofReqLabel] = presID

	return nil
}

func (co *InstitutionController) DID() ssi.DIDPrism {
	return co.InstitutionDIDPrism
}

func (co *InstitutionController) ConnectionsHashMap() map[string]ssi.ConnectionID {
	return co.Connections
}

func (co *InstitutionController) CredentialOffersReceivedSlice() []ssi.RecordID {
	return co.CredentialOffersReceived
}

func (co *InstitutionController) CredentialsHashMap() map[string]ssi.RecordID {
	return co.Credentials
}

func (co *InstitutionController) ProofRequestsAcceptedHashMap() map[string]ssi.PresentationID {
	return co.ProofRequestsAccepted
}

func (co *InstitutionController) ProofRequestsReceivedHashMap() []ssi.PresentationID {
	return co.ProofRequestsReceived
}
