package domain

import "github.com/zafir0101/SSI-ENV/internal/ssi"

type Controller interface {
	DID() string

	RefreshOffersReceived() error
	CredentialsHashMap() map[string]ssi.RecordID
	CredentialOffersReceivedSlice() []ssi.RecordID
	AcceptCredentialOffer(credentialLabel string, recID ssi.RecordID) error

	RefreshProofRequests() error
	AcceptProofRequest(proofReqLabel string, credLabel string, presID ssi.PresentationID) error
	ProofRequestsAcceptedHashMap() map[string]ssi.PresentationID
	ProofRequestsReceivedHashMap() []ssi.PresentationID

	CreateConnection(connectionLabel string) (ssi.InvitationOOB, error)
	AcceptConnection(connectionLabel string, invOOB ssi.InvitationOOB) error
	ConnectionsHashMap() map[string]ssi.ConnectionID
}
