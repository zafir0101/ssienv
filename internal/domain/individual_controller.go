package domain

import (
	"errors"

	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

type IndividualController struct {
	edgeAgentAPI *ssi.EdgeAgentAPI

	individualDIDPrism ssi.LongFormDIDPrism

	connections map[string]ssi.ConnectionID // Será serializado, apenas connections realizados na maquina

	credentials              map[string]ssi.RecordID
	credentialOffersReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a oferta

	proofRequestsAccepted map[string]ssi.PresentationID
	proofRequestsReceived []ssi.RecordID // Limitação: Não consegue compartilhar uma label para a requisicao
}

func NewIndividualController(cloudAgentAPI *ssi.EdgeAgentAPI) *IndividualController {
	return &IndividualController{
		edgeAgentAPI:          cloudAgentAPI,
		connections:           make(map[string]ssi.ConnectionID),
		credentials:           make(map[string]ssi.RecordID),
		proofRequestsAccepted: make(map[string]ssi.PresentationID),
	}
}

func (co *IndividualController) RefreshOffersReceived() error {
	recordIDs, credentialStatus, err := co.edgeAgentAPI.ListCredentialOffers()
	if err != nil {
		return err
	}

	clear(co.credentialOffersReceived)
	co.credentialOffersReceived = co.credentialOffersReceived[:0]

	for i, credStatus := range credentialStatus {
		if credStatus == "OfferReceived" {
			co.credentialOffersReceived = append(co.credentialOffersReceived, recordIDs[i])
		}
	}

	return nil
}

func (co *IndividualController) RefreshProofRequestsReceived() error {
	presentationIDs, proofReqStatus, err := co.edgeAgentAPI.ListProofRequestsData()
	if err != nil {
		return err
	}

	clear(co.proofRequestsReceived)
	co.proofRequestsReceived = co.proofRequestsReceived[:0]

	for i, proofReqStatus := range proofReqStatus {
		if proofReqStatus == "RequestReceived" {
			co.proofRequestsReceived = append(co.proofRequestsReceived, presentationIDs[i])
			continue
		}
	}

	return nil
}

func (co *IndividualController) CreateDID() error {
	pksID := []string{"key1-authentication", "key2-assertionMethod"}
	pksPurpose := []KeyPurpose{Authentication, AssertionMethod}

	payload, err := newDIDCreationPayload(pksID, pksPurpose)
	if err != nil {
		return err
	}

	did, err := co.edgeAgentAPI.CreateDID(payload)
	if err != nil {
		return err
	}

	co.individualDIDPrism = did

	return nil
}

func (co *IndividualController) DeactivateDID() error {
	if co.individualDIDPrism == "" {
		return errors.New("First create a did")
	}

	if err := co.edgeAgentAPI.DeactivateDID(co.individualDIDPrism); err != nil {
		return err
	}

	return nil
}

func (co *IndividualController) CreateConnection(label string) (ssi.InvitationOOB, error) {
	payload := newConnectionCreationPayload(label)

	connID, invOOB, err := co.edgeAgentAPI.CreateConnection(payload)
	if err != nil {
		return "", err
	}

	co.connections[label] = connID
	return invOOB, nil
}

func (co *IndividualController) AcceptConnection(label string, invOOB ssi.InvitationOOB) error {
	payload := newConnectionAcceptPayload(invOOB)

	connID, err := co.edgeAgentAPI.AcceptConnection(payload)
	if err != nil {
		return err
	}

	co.connections[label] = connID
	return nil
}

func (co *IndividualController) DeactivateConnection(label string) error {
	connID := co.connections[label]
	if connID == "" {
		return errors.New("No connections with label " + label)
	}

	if err := co.edgeAgentAPI.DeactivateConnection(connID); err != nil {
		return err
	}

	delete(co.connections, label)

	return nil
}

func (co *IndividualController) AcceptCredentialOffer(credentialLabel string, recID ssi.RecordID) error {
	payload := newOfferAcceptancePayload(co.individualDIDPrism)

	if err := co.edgeAgentAPI.AcceptCredentialOffer(payload, recID); err != nil {
		return err
	}

	for i, offer := range co.credentialOffersReceived {
		if offer == recID {
			co.credentialOffersReceived = append(co.credentialOffersReceived[:i], co.credentialOffersReceived[i+1:]...)
			break
		}
	}

	co.credentials[credentialLabel] = recID

	return nil
}

func (co *IndividualController) AcceptProofRequest(proofReqLabel string, credentialLabel string, presID ssi.PresentationID) error {
	recID := co.credentials[credentialLabel]
	if recID == "" {
		return errors.New("No credentials with label " + credentialLabel)
	}

	payload := newProofRequestAcceptancePayload(recID)

	if err := co.edgeAgentAPI.AcceptProofRequest(payload, presID); err != nil {
		return err
	}

	for i, request := range co.proofRequestsReceived {
		if request == presID {
			co.proofRequestsReceived = append(co.proofRequestsReceived[:i], co.proofRequestsReceived[i+1:]...)
			break
		}
	}

	co.proofRequestsAccepted[proofReqLabel] = presID

	return nil
}

// Getters
func (co *IndividualController) Connections() map[string]ssi.ConnectionID {
	return co.connections
}

func (co *IndividualController) CrendetialOffersReceived() []ssi.RecordID {
	return co.credentialOffersReceived
}

func (co *IndividualController) Credentials() map[string]ssi.RecordID {
	return co.credentials
}

func (co *IndividualController) ProofRequestsReceived() []ssi.PresentationID {
	return co.proofRequestsReceived
}
