package ssi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CloudAgentAPI struct {
	AgentURL *url.URL
}

func NewCloudAgentAPI(agentURL *url.URL) *CloudAgentAPI {
	return &CloudAgentAPI{
		AgentURL: agentURL,
	}
}

func (ca *CloudAgentAPI) CreateDID(payload Payload) (DIDPrism, error) {
	longFormDID, err := registerDID(payload, formatURL(ca.AgentURL))
	if err != nil {
		return "", err
	}

	did, err := ca.PublishDID(longFormDID)
	if err != nil {
		return "", err
	}

	return did, nil
}

func (ca *CloudAgentAPI) PublishDID(longFormDID LongFormDIDPrism) (DIDPrism, error) {
	respPub, err := http.Post(formatURL(ca.AgentURL)+"/did-registrar/dids/"+longFormDID+"/publications",
		"application/json", nil)
	if err != nil {
		return "", err
	}
	defer respPub.Body.Close()

	if respPub.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(respPub.Body)
		return "", errors.New("publishing failed: status=" + respPub.Status + "body=" + string(body))
	}

	var didPubResponse didPubResponse
	if err := json.NewDecoder(respPub.Body).Decode(&didPubResponse); err != nil {
		return "", err
	}
	did := DIDPrism(didPubResponse.ScheduledOperation.DIDRef)

	return did, nil
}

func (ca *CloudAgentAPI) ResolveDID(did DIDPrism) (DIDPrismDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, formatURL(ca.AgentURL)+"/dids/"+did, nil)
	if err != nil {
		return DIDPrismDocument{}, err
	}

	req.Header.Set("Accept", "application/ld+json; profile=https://w3id.org/did-resolution")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DIDPrismDocument{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return DIDPrismDocument{}, errors.New("resolution failed: status=" + resp.Status + "body=" + string(body))
	}

	var didDocument DIDPrismDocument
	if err := json.NewDecoder(resp.Body).Decode(&didDocument); err != nil {
		return DIDPrismDocument{}, err
	}

	return didDocument, nil
}

// Limitada em adicionar ou remover chaves
func (ca *CloudAgentAPI) UpdateDID(payload Payload, did DIDPrism) error {
	postBody, err := toIOReader(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(formatURL(ca.AgentURL)+"/did-registrar/dids/"+did+"/updates",
		"application/json", postBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("update failed: status=" + resp.Status + "body=" + string(body))
	}

	return nil
}

// Retorna 202 mas não é efetivada na VDR no ambiente de teste locais
// Verificar se desativa quando não publicado
func (ca *CloudAgentAPI) DeactivateDID(did DIDPrism) error {
	resp, err := http.Post(formatURL(ca.AgentURL)+"/did-registrar/dids/"+did+"/deactivations",
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

func (ca *CloudAgentAPI) CreateConnection(payload Payload) (ConnectionID, InvitationOOB, error) {
	return createConnection(payload, formatURL(ca.AgentURL))
}

func (ca *CloudAgentAPI) AcceptConnection(payload Payload) (ConnectionID, error) {
	return acceptConnection(payload, formatURL(ca.AgentURL))
}

// Limitado a convites enviados mas não respondidos.
func (ca *CloudAgentAPI) DeactivateConnection(connID ConnectionID) error {
	return deactivateConnection(connID, formatURL(ca.AgentURL))
}

func (ca *CloudAgentAPI) CreateSchema(payload Payload) (SchemaID, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(formatURL(ca.AgentURL)+"/schema-registry/schemas",
		"application/json", postBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("schema creation failed: status=" + resp.Status + "body=" + string(body))
	}

	var schemaResp schemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&schemaResp); err != nil {
		return "", err
	}

	schemaID := formatURL(ca.AgentURL) + "/schema-registry/schemas/" + schemaResp.SchemaGUID + "/schema"
	return schemaID, nil
}

func (ca *CloudAgentAPI) CreateCredentialOffer(payload Payload) (RecordID, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(formatURL(ca.AgentURL)+"/issue-credentials/credential-offers",
		"application/json", postBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("credendial offer creation failed: status=" + resp.Status + "body=" + string(body))
	}

	var content credentialOfferContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return "", err
	}

	return content.RecordID, nil
}

func (ca *CloudAgentAPI) ListCredentialOffers() ([]RecordID, []CredentialStatus, error) {
	return listCredentialOffers(formatURL(ca.AgentURL))
}

func (ca *CloudAgentAPI) AcceptCredentialOffer(payload Payload, recID RecordID) error {
	return acceptCredentialOffer(payload, recID, formatURL(ca.AgentURL))
}

func (ca *CloudAgentAPI) CreateProofRequest(payload Payload) (PresentationID, error) {
	postBody, err := toIOReader(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(formatURL(ca.AgentURL)+"/present-proof/presentations",
		"application/json", postBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("proof request creation failed: status=" + resp.Status + "body=" + string(body))
	}

	var proofReqRes proofRequestContent
	if err := json.NewDecoder(resp.Body).Decode(&proofReqRes); err != nil {
		return "", err
	}

	return proofReqRes.PresentationID, nil
}

func (ca *CloudAgentAPI) ListProofRequestsData() ([]PresentationID, []ProofRequestStatus, error) {
	return listProofRequestsData(formatURL(ca.AgentURL))
}

func (ca *CloudAgentAPI) AcceptProofRequest(payload Payload, presID PresentationID) error {
	return acceptProofRequest(payload, presID, formatURL(ca.AgentURL))
}
