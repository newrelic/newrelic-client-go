package aimonitoring

import "context"

// DO NOT DELETE
//
// This file holds code that tutone cannot generate correctly for this package.
// Regenerating pkg/aimonitoring should leave this file intact.

// AimDatasetUpsertRecordsInput - The input for upserting records into a dataset, which creates a new version.
//
// Hand-written (skip_type_create in .tutone.yml): tutone generates Insert as a struct value,
// which encoding/json never omits, so every upsert would send `insert: {}` and fail validation.
type AimDatasetUpsertRecordsInput struct {
	// The identifier of the dataset to upsert records into.
	DatasetId string `json:"datasetId"`
	// The identifiers of existing records to delete.
	DeleteRecordIds []int `json:"deleteRecordIds,omitempty"`
	// Records to insert at a specific position, prepend to the beginning of the dataset if afterID is 0, or appended to the end if afterId is omitted.
	Insert *AimDatasetInsertRecordInput `json:"insert,omitempty"`
	// An optional message describing the new version.
	Message string `json:"message,omitempty"`
	// Records to replace the dataset's existing records with.
	Records []AimDatasetRecordInput `json:"records,omitempty"`
}

// AimDatasetInsertRecordInput - Records to insert into a dataset at a specific position.
//
// Hand-written (skip_type_create in .tutone.yml): AfterId is a pointer because 0 (prepend)
// and omitted (append) mean different things, which `int,omitempty` cannot express.
type AimDatasetInsertRecordInput struct {
	// The identifier of the record after which the new records will be inserted.
	AfterId *int `json:"afterId,omitempty"`
	// The records to insert.
	Records []AimDatasetRecordInput `json:"records"`
}

// Fetches a page of records from a specific dataset version.
//
// Hand-written: the query lives under actor.account, and tutone cannot combine that path
// with actor.entityManagement in one package without emitting a duplicate Actor type.
func (a *Aimonitoring) GetRecords(
	accountID int,
	aimDatasetId string,
	versionId string,
	cursor *string,
	recordIds []int,
) (*AimDatasetRecordsResult, error) {
	return a.GetRecordsWithContext(context.Background(),
		accountID,
		aimDatasetId,
		versionId,
		cursor,
		recordIds,
	)
}

// Fetches a page of records from a specific dataset version.
func (a *Aimonitoring) GetRecordsWithContext(
	ctx context.Context,
	accountID int,
	aimDatasetId string,
	versionId string,
	cursor *string,
	recordIds []int,
) (*AimDatasetRecordsResult, error) {

	resp := recordsResponse{}
	vars := map[string]interface{}{
		"accountID":    accountID,
		"aimDatasetId": aimDatasetId,
		"cursor":       cursor,
		"recordIds":    recordIds,
		"versionId":    versionId,
	}

	if err := a.client.NerdGraphQueryWithContext(ctx, getRecordsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.AimDataset.Records, nil
}

type recordsResponse struct {
	Actor struct {
		Account struct {
			AimDataset AimDatasetAccountStitchedFields `json:"aimDataset"`
		} `json:"account"`
	} `json:"actor"`
}

// AimDatasetAccountStitchedFields - Account-scoped fields for the AimDataset API.
type AimDatasetAccountStitchedFields struct {
	// Fetches a page of records from a specific dataset version.
	Records AimDatasetRecordsResult `json:"records,omitempty"`
}

const getRecordsQuery = `query(
	$accountID: Int!,
	$aimDatasetId: ID!,
	$cursor: String,
	$recordIds: [Int!],
	$versionId: ID!,
) { actor { account(id: $accountID) { aimDataset { records(
	aimDatasetId: $aimDatasetId,
	cursor: $cursor,
	recordIds: $recordIds,
	versionId: $versionId,
) {
	nextCursor
	records {
		additionalFields
		context
		expected
		id
		input
	}
	totalCount
} } } } }`
