package domain

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type InstitutionController struct {
	CloudAgentAPI *ssi.CloudAgentAPI

	InstitutionDIDPrism ssi.DIDPrism
	PublishedsDIDs      map[string]ssi.DIDPrism // Sera serializado, apenas DIDs publicados na máquina

	Connections map[string]ssi.ConnectionID // Será serializado, apenas connections realizados na maquina

	Shemas map[string]ssi.SchemaID // Será serializado, apenas schemas criados na maquina

	Credentials              map[string]ssi.RecordID
	CredentialOffersReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a oferta
	CredentialOffersSent     map[string]ssi.RecordID

	ProofRequestsAccepted           map[string]ssi.PresentationID
	ProofRequestsAcceptedByExternal []ssi.RecordID
	ProofRequestsReceived           []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a requisicao
	ProofRequestsSent               map[string]ssi.PresentationID

	Num_keys int
}

func NewInstitutionController(cloudAgentAPI *ssi.CloudAgentAPI) *InstitutionController {
	return &InstitutionController{
		CloudAgentAPI:         cloudAgentAPI,
		PublishedsDIDs:        make(map[string]ssi.DIDPrism),
		Connections:           make(map[string]ssi.ConnectionID),
		Shemas:                make(map[string]ssi.SchemaID),
		Credentials:           make(map[string]ssi.RecordID),
		CredentialOffersSent:  make(map[string]ssi.RecordID),
		ProofRequestsAccepted: make(map[string]ssi.PresentationID),
		ProofRequestsSent:     make(map[string]ssi.PresentationID),
	}
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

func (co *InstitutionController) RefreshProofRequestsReceived() error {
	presentationIDs, proofReqStatus, err := co.CloudAgentAPI.ListProofRequestsData()
	if err != nil {
		return err
	}

	clear(co.ProofRequestsReceived)
	co.ProofRequestsReceived = co.ProofRequestsReceived[:0]
	clear(co.ProofRequestsAcceptedByExternal)
	co.ProofRequestsAcceptedByExternal = co.ProofRequestsAcceptedByExternal[:0]

	for i, proofReqStatus := range proofReqStatus {
		if proofReqStatus == "RequestReceived" {
			co.ProofRequestsReceived = append(co.ProofRequestsReceived, presentationIDs[i])
			continue
		}
		if proofReqStatus == "PresentationVerified" {
			co.ProofRequestsAcceptedByExternal = append(co.ProofRequestsAcceptedByExternal, presentationIDs[i])
		}
	}

	return nil
}

func (co *InstitutionController) CreateDID() error {
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

	co.PublishedsDIDs[label] = did
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
	if co.InstitutionDIDPrism == "" {
		return errors.New("First create a did")
	}

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

	return nil
}

func (co *InstitutionController) RemoveDIDKey(pkID string, pkPurpose KeyPurpose) error {
	if co.InstitutionDIDPrism == "" {
		return errors.New("First create a did")
	}

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

	return nil
}

func (co *InstitutionController) DeactivateDID() error {
	if co.InstitutionDIDPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.CloudAgentAPI.DeactivateDID(co.InstitutionDIDPrism); err != nil {
		return err
	}

	return nil
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

func (co *InstitutionController) DeactivateConnection(label string) error {
	connID := co.Connections[label]
	if connID == "" {
		return errors.New("No connections with label " + label)
	}

	if err := co.CloudAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	delete(co.Connections, label)

	return nil
}

func (co *InstitutionController) CreateSchema(schemaName string, schema json.RawMessage) error {
	payload := newSchemaCreationPayload(schemaName, co.InstitutionDIDPrism, schema)

	schemaID, err := co.CloudAgentAPI.CreateSchema(payload)
	if err != nil {
		return err
	}

	co.Shemas[schemaName] = schemaID

	return nil
}

func (co *InstitutionController) CreateCredentialOffer(offerLabel string, claims json.RawMessage,
	connLabel string, schemaID ssi.SchemaID) error {
	connID := co.Connections[connLabel]
	if connID == "" {
		return errors.New("No connections with label " + connLabel)
	}

	payload := newCredentialOfferPayload(claims, co.InstitutionDIDPrism, connID, schemaID)

	recordID, err := co.CloudAgentAPI.CreateCredentialOffer(payload)
	if err != nil {
		return err
	}

	co.CredentialOffersSent[offerLabel] = recordID

	return nil
}

func (co *InstitutionController) AcceptCredentialOffer(credentialLabel string, recID ssi.RecordID) error {
	payload := newOfferAcceptancePayload(co.InstitutionDIDPrism)

	if err := co.CloudAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	for i, offer := range co.CredentialOffersReceived {
		if offer == recID {
			co.CredentialOffersReceived = append(co.CredentialOffersReceived[:i], co.CredentialOffersReceived[i+1:]...)
			break
		}
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
		return errors.New("No connections with label " + connLabel)
	}

	payload := newProofRequestPayload(proofReqLabel, connID, schemaID, co.InstitutionDIDPrism)

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

	payload := newProofRequestAcceptancePayload(recID)

	if err := co.CloudAgentAPI.AcceptProofRequest(payload, presID); err != nil {
		return err
	}

	for i, request := range co.ProofRequestsReceived {
		if request == presID {
			co.ProofRequestsReceived = append(co.ProofRequestsReceived[:i], co.ProofRequestsReceived[i+1:]...)
			break
		}
	}

	co.ProofRequestsAccepted[proofReqLabel] = presID

	return nil
}
