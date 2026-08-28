package ssi

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type EdgeAgentAPI struct {
	AgentURL *url.URL
}

func NewEdgeAgentAPI(agentURL *url.URL) *EdgeAgentAPI {
	return &EdgeAgentAPI{
		AgentURL: agentURL,
	}
}

func (ea *EdgeAgentAPI) CreateDID(payload Payload) (LongFormDIDPrism, error) {
	return registerDID(payload, formatURL(ea.AgentURL))
}

// Retorna 202 mas não é efetivada na VDR no ambiente de teste locais
// Verificar se desativa quando não publicado
func (ea *EdgeAgentAPI) DeactivateDID(did DIDPrism) error {
	resp, err := http.Post(formatURL(ea.AgentURL)+"/did-registrar/dids/"+did+"/deactivations",
		"application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("deactivation failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

func (ea *EdgeAgentAPI) CreateConnection(payload Payload) (ConnectionID, InvitationOOB, error) {
	return createConnection(payload, formatURL(ea.AgentURL))
}

func (ea *EdgeAgentAPI) AcceptConnection(payload Payload) (ConnectionID, error) {
	return acceptConnection(payload, formatURL(ea.AgentURL))
}

// Limitado a convites enviados mas não respondidos.
func (ea *EdgeAgentAPI) DeactivateConnection(connID ConnectionID) error {
	return deactivateConnection(connID, formatURL(ea.AgentURL))
}

func (ea *EdgeAgentAPI) AcceptProofRequest(payload Payload, presID PresentationID) error {
	return acceptProofRequest(payload, presID, formatURL(ea.AgentURL))
}

func (ea *EdgeAgentAPI) ListProofRequestsData() ([]PresentationID, []ProofRequestStatus, error) {
	return listProofRequestsData(formatURL(ea.AgentURL))
}

func (ea *EdgeAgentAPI) AcceptCredentialOffer(payload Payload, recID RecordID) error {
	return acceptCredentialOffer(payload, recID, formatURL(ea.AgentURL))
}

func (ea *EdgeAgentAPI) ListCredentialOffers() ([]RecordID, []CredentialStatus, error) {
	return listCredentialOffers(formatURL(ea.AgentURL))
}
