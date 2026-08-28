package domain

import (
	"errors"

	"github.com/zafir0101/ssienv/internal/ssi"
)

type IndividualController struct {
	EdgeAgentAPI *ssi.EdgeAgentAPI

	IndividualDIDPrism ssi.LongFormDIDPrism

	Connections map[string]ssi.ConnectionID // Será serializado, apenas connections realizados na maquina

	Credentials              map[string]ssi.RecordID
	CredentialOffersReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a oferta

	ProofRequestsAccepted map[string]ssi.PresentationID
	ProofRequestsReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a requisicao
}

func NewIndividualController(cloudAgentAPI *ssi.EdgeAgentAPI) (*IndividualController, error) {
	controller := &IndividualController{
		EdgeAgentAPI:          cloudAgentAPI,
		Connections:           make(map[string]ssi.ConnectionID),
		Credentials:           make(map[string]ssi.RecordID),
		ProofRequestsAccepted: make(map[string]ssi.PresentationID),
	}
	if err := controller.createDID(); err != nil {
		return nil, err
	}

	return controller, nil
}

func (co *IndividualController) RefreshOffersReceived() error {
	recordIDs, credentialStatus, err := co.EdgeAgentAPI.ListCredentialOffers()
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

func (co *IndividualController) RefreshProofRequests() error {
	presentationIDs, proofReqStatus, err := co.EdgeAgentAPI.ListProofRequestsData()
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
	}

	return nil
}

func (co *IndividualController) createDID() error {
	pksID := []string{"key1-authentication", "key2-assertionMethod"}
	pksPurpose := []KeyPurpose{Authentication, AssertionMethod}

	payload, err := newDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return err
	}

	did, err := co.EdgeAgentAPI.CreateDID(payload)
	if err != nil {
		return err
	}

	co.IndividualDIDPrism = did

	return nil
}

func (co *IndividualController) DeactivateDID() error {
	if co.IndividualDIDPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.EdgeAgentAPI.DeactivateDID(co.IndividualDIDPrism); err != nil {
		return err
	}

	return nil
}

func (co *IndividualController) CreateConnection(label string) (ssi.InvitationOOB, error) {
	payload := newConnectionCreationPayload(label)

	connID, invOOB, err := co.EdgeAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", err
	}

	co.Connections[label] = connID
	return invOOB, nil
}

func (co *IndividualController) AcceptConnection(label string, invOOB ssi.InvitationOOB) error {
	payload := newConnectionAcceptPayload(invOOB)

	connID, err := co.EdgeAgentAPI.AcceptConnection(payload)
	if err != nil {
		return err
	}

	co.Connections[label] = connID
	return nil
}

func (co *IndividualController) DeactivateConnection(label string) error {
	connID := co.Connections[label]
	if connID == "" {
		return errors.New("No connections with label " + label)
	}

	if err := co.EdgeAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	delete(co.Connections, label)

	return nil
}

func (co *IndividualController) AcceptCredentialOffer(credentialLabel string, recID ssi.RecordID) error {
	payload := newOfferAcceptancePayload(co.IndividualDIDPrism)

	if err := co.EdgeAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
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

func (co *IndividualController) AcceptProofRequest(proofReqLabel string, credentialLabel string, presID ssi.PresentationID) error {
	recID := co.Credentials[credentialLabel]
	if recID == "" {
		return errors.New("No credentials with label " + credentialLabel)
	}

	payload := newProofRequestAcceptancePayload(recID)

	if err := co.EdgeAgentAPI.AcceptProofRequest(payload, presID); err != nil {
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

func (co *IndividualController) DID() ssi.LongFormDIDPrism {
	return co.IndividualDIDPrism
}

func (co *IndividualController) ConnectionsHashMap() map[string]ssi.ConnectionID {
	return co.Connections
}

func (co *IndividualController) CredentialOffersReceivedSlice() []ssi.RecordID {
	return co.CredentialOffersReceived
}

func (co *IndividualController) CredentialsHashMap() map[string]ssi.RecordID {
	return co.Credentials
}

func (co *IndividualController) ProofRequestsAcceptedHashMap() map[string]ssi.PresentationID {
	return co.ProofRequestsAccepted
}

func (co *IndividualController) ProofRequestsReceivedHashMap() []ssi.PresentationID {
	return co.ProofRequestsReceived
}
