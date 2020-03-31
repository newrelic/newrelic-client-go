package nerdstorage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

type EntityScopedDoc struct {
	MyField string
}

func Example_entityScope() {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := New(cfg)

	entityID := "MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
	packageID := "ecaeb28c-7b3f-4932-9e33-7385980efa1c"

	// Write a NerdStorage document with account scope.
	writeDocumentInput := WriteDocumentInput{
		PackageID:  packageID,
		Collection: "myCol",
		DocumentID: "myDoc",
		Document: AccountScopedDoc{
			MyField: "myValue",
		},
	}

	_, err := client.WriteDocumentWithEntityScope(entityID, writeDocumentInput)
	if err != nil {
		log.Fatal("error writing NerdStorage document:", err)
	}

	// Write a second NerdStorage document to the same collection with account scope.
	writeAlternateDocumentInput := writeDocumentInput
	writeAlternateDocumentInput.DocumentID = "myOtherDoc"

	_, err = client.WriteDocumentWithEntityScope(entityID, writeAlternateDocumentInput)
	if err != nil {
		log.Fatal("error writing NerdStorage document:", err)
	}

	// Get a NerdStorage collection with account scope.
	getCollectionInput := GetCollectionInput{
		PackageID:  packageID,
		Collection: "myCol",
	}

	collection, err := client.GetCollectionWithEntityScope(entityID, getCollectionInput)
	if err != nil {
		log.Fatal("error retrieving NerdStorage collection:", err)
	}

	fmt.Printf("Collection length: %v\n", len(collection))

	// Get a NerdStorage document with account scope.
	getDocumentInput := GetDocumentInput{
		PackageID:  packageID,
		Collection: "myCol",
		DocumentID: "myDoc",
	}

	rawDoc, err := client.GetDocumentWithEntityScope(entityID, getDocumentInput)
	if err != nil {
		log.Fatal("error retrieving NerdStorage document:", err)
	}

	// Convert the document to a struct.
	var myDoc AccountScopedDoc

	// Method 1:
	marshalled, err := json.Marshal(rawDoc)
	if err != nil {
		log.Fatal("error marshalling NerdStorage document to json:", err)
	}

	err = json.Unmarshal(marshalled, &myDoc)
	if err != nil {
		log.Fatal("error unmarshalling NerdStorage document to struct:", err)
	}

	fmt.Printf("Document: %v\n", myDoc)

	// Method 2:
	err = mapstructure.Decode(rawDoc, &myDoc)
	if err != nil {
		log.Fatal("error converting NerdStorage document to struct:", err)
	}

	fmt.Printf("Document: %v\n", myDoc)

	// Delete a NerdStorage document with account scope.
	deleteDocumentInput := DeleteDocumentInput{
		PackageID:  packageID,
		Collection: "myCol",
		DocumentID: "myDoc",
	}

	ok, err := client.DeleteDocumentWithEntityScope(entityID, deleteDocumentInput)

	if !ok || err != nil {
		log.Fatal("error deleting NerdStorage document:", err)
	}

	// Delete a NerdStorage collection with account scope.
	deleteCollectionInput := DeleteCollectionInput{
		PackageID:  packageID,
		Collection: "myCol",
	}

	ok, err = client.DeleteCollectionWithEntityScope(entityID, deleteCollectionInput)
	if err != nil {
		log.Fatal("error deleting NerdStorage collection:", err)
	}
}
