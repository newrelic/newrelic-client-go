package aimonitoring

// DO NOT DELETE
//
// This file holds code that tutone cannot generate correctly for this package.
// Regenerating pkg/aimonitoring should leave this file intact.

// AimDatasetInsertRecordInput - Records to insert into a dataset at a specific position.
//
// Hand-written (skip_type_create in .tutone.yml): AfterId is a pointer because 0 (prepend)
// and omitted (append) mean different things, which `int,omitempty` cannot express.
type AimDatasetInsertRecordInput struct {
	// The identifier of the record after which the new records will be inserted.
	// Use 0 to prepend to the beginning of the dataset.
	// Omit or pass null to append to the end of the dataset.
	AfterId *int `json:"afterId,omitempty"`
	// The records to insert.
	Records []AimDatasetRecordInput `json:"records"`
}
