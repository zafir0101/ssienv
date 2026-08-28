package ssi_test

/*
import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/zafir0101/SSI-ENV/internal/domain"
	"github.com/zafir0101/SSI-ENV/internal/ssi"
)

const caUrlString = "http://host.docker.internal:8080"
const eaUrlString = "http://host.docker.internal:8081"

var agentURL, _ = url.Parse(eaUrlString)

var ca *ssi.CloudAgentAPI = createAgent(caUrlString)
var ea *ssi.EdgeAgentAPI = ssi.NewEdgeAgentAPI(agentURL)
var ins_co *domain.InstitutionController = domain.NewInstitutionController(ca)
var ind_co *domain.IndividualController = domain.NewIndividualController(ea)

func createAgent(urlString string) *ssi.CloudAgentAPI {
	agentURL, _ := url.Parse(urlString)
	return ssi.NewCloudAgentAPI(agentURL)
}

func TestCrudDID(t *testing.T) {
	// Teste de Create DID
	err := ins_co.CreateDID()
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(30 * time.Second)

	// Teste de Resolve DID
	/*
		didDocument, err := co.ResolveDID("did:prism:318039f23c7c4356b8650a42a05a72b7f1605b8689accac3a995363d5d80645b")
		if err != nil {
			fmt.Println(err.Error())
		}

		json, err := json.MarshalIndent(didDocument, "", "  ")
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(string(json))
*/

// Teste de Deactivate DID
/*
	if err := ca.DeactivateDID("did:prism:318039f23c7c4356b8650a42a05a72b7f1605b8689accac3a995363d5d80645b"); err != nil {
		fmt.Println(err.Error())
	}
*/
/*
	   // Teste de Update DID
	   		if err := co.AddKeyToDID(0); err != nil {
	   			fmt.Println(err.Error())
	   		}

	   		time.Sleep(30 * time.Second)

	   		if err := co.RemoveDIDKey("key3-authentication", 0); err != nil {
	   			fmt.Println(err.Error())
	   		}


	   	// Teste de Connections
	   	connLabel := "ConectarNoVinizao"
	   	inv, err := ins_co.CreateConnection(connLabel)
	   	if err != nil {
	   		fmt.Println(err.Error())
	   	}

	   	time.Sleep(15 * time.Second)

	   	ind_co.AcceptConnection(connLabel, inv)

	   	// if err := co.DeactivateConnection(label); err != nil {
	   	// fmt.Println(err.Error())
	   	// }

	   	// Teste para Schemas
	   	schema := json.RawMessage(`
	   			    {
	   			        "$id": "https://example.com/driving-license-1.0",
	   			        "$schema": "https://json-schema.org/draft/2020-12/schema",
	   			        "type": "object",
	   			        "properties": {
	   			            "emailAddress": {"type": "string"},
	   			            "givenName": {"type": "string"},
	   			            "familyName": {"type": "string"}
	   			        },
	   			        "required": ["emailAddress", "givenName", "familyName"],
	   			        "additionalProperties": false
	   			    }`)

	   	label := "Credencial de Cortesã do vini"
	   	if err := ins_co.CreateSchema(label, schema); err != nil {
	   		fmt.Println(err.Error())
	   	}

	   	// fmt.Println(co.RetrieveSchemas()[label])

	   	time.Sleep(15 * time.Second)

	   	// Teste para CreateCredentialOffer
	   	claims := json.RawMessage(`{
	   		    "emailAddress": "zeca.galhao@gmail.com",
	   		    "givenName": "Zeca Galhao",
	   		    "familyName": "galhao"
	   		  }`)
	   	offerlabel := "Para provar a pureza do vini"
	   	if err := ins_co.CreateCredentialOffer(offerlabel, claims, connLabel, ins_co.Schemas()[label]); err != nil {
	   		fmt.Println(err.Error())
	   	}

	   	time.Sleep(15 * time.Second)

	   	// Teste para receber oferta de credencial
	   	// O edge agent vai utilizar o longform (nao ira publicar o did)
	   	ind_co.CreateDID()

	   	time.Sleep(30 * time.Second)

	   	var recIDs []string
	   	for i := 0; i < 10; i++ {
	   		err := ind_co.RefreshOffersReceived()
	   		if err != nil && err.Error() != "No credential offers" {
	   			t.Fatalf("refresh offers failed: %v", err) // erro de verdade, para
	   		}
	   		recIDs = ind_co.CrendetialOffersReceived()
	   		if len(recIDs) > 0 {
	   			break
	   		}
	   		time.Sleep(3 * time.Second)
	   	}

	   	credlabel := "Credencial que vini é puro"
	   	if err := ind_co.AcceptCredentialOffer(credlabel, recIDs[0]); err != nil {
	   		fmt.Println(err.Error())
	   	}

	   	// Teste para apresentar prova de credencial
	   	prooflabel := "Provar a pureza do vini"
	   	err = ins_co.CreateProofRequest(prooflabel, connLabel, ins_co.Schemas()[label])
	   	if err != nil {
	   		fmt.Println(err.Error())
	   	}

	   	// Teste para aceitar a apresentação de prova de credencial
	   	var presIDs []string
	   	for i := 0; i < 10; i++ {
	   		if err := ind_co.RefreshProofRequestsReceived(); err != nil {
	   			t.Fatalf("refresh proof requests failed: %v", err)
	   		}
	   		presIDs = ind_co.ProofRequestsReceived()
	   		if len(presIDs) > 0 {
	   			break
	   		}
	   		time.Sleep(3 * time.Second)
	   	}
	   	if len(presIDs) == 0 {
	   		t.Fatalf("nenhuma proof request recebida após 30s")
	   	}
	   	err = ind_co.AcceptProofRequest(prooflabel, credlabel, presIDs[0])
	   	if err != nil {
	   		fmt.Println(err.Error())
	   	}
}
*/
