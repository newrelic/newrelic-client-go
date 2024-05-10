//go:build integration
// +build integration

package entities

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestDecodeEntityGuid_Valid(t *testing.T) {
	entity := GenericEntity{
		AccountId:  12345,
		Domain:     "test_domain",
		EntityType: "user",
		DomainId:   "abc123",
	}
	encodedGuid := encodeEntity(entity)

	decodedEntity, err := DecodeEntityGuid(encodedGuid)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if decodedEntity.AccountId != entity.AccountId ||
		decodedEntity.Domain != entity.Domain ||
		decodedEntity.EntityType != entity.EntityType ||
		decodedEntity.DomainId != entity.DomainId {
		t.Errorf("Decoded entity does not match original entity")
	}
}

func encodeEntity(entity GenericEntity) string {
	parts := []string{
		strconv.FormatInt(entity.AccountId, 10),
		entity.Domain,
		entity.EntityType,
		entity.DomainId,
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(parts, DELIMITER)))
}

func TestDecodeEntityGuid_MissingDelimiter(t *testing.T) {
	invalidGuid := "invalidentityguid"

	_, err := DecodeEntityGuid(invalidGuid)

	if err != EntityGUIDValidationErrorTypes.INVALID_ENTITY_GUID_ERROR {
		t.Errorf("Expected error 'invalid entity GUID format', got %v", err)
	}
}

func TestDecodeEntityGuid_LessThanFourParts(t *testing.T) {
	invalidGuid := base64.StdEncoding.EncodeToString([]byte("account|domain"))

	_, err := DecodeEntityGuid(invalidGuid)

	expectedErrorMessage := fmt.Sprintf("invalid entity GUID format: expected at least 4 parts delimited by '%s': %s", DELIMITER, invalidGuid)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message: %s, got %v", expectedErrorMessage, err)
	}
}

func TestDecodeEntityGuid_EmptyEntityType(t *testing.T) {
	encodedGuid := base64.StdEncoding.EncodeToString([]byte("12345|domain||domainId"))

	_, err := DecodeEntityGuid(encodedGuid)

	if err != EntityGUIDValidationErrorTypes.EMPTY_ENTITY_TYPE_ERROR {
		t.Errorf("Expected error 'entity type is required', got %v", err)
	}
}
