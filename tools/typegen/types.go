package main

type Kind string

const (
	KindENUM        Kind = "ENUM"
	KindInputObject Kind = "INPUT_OBJECT"
	KindInterface   Kind = "INTERFACE"
	KindList        Kind = "LIST"
	KindNonNull     Kind = "NON_NULL"
	KindObject      Kind = "OBJECT"
	KindScalar      Kind = "SCALAR"
)

// Schema contains data about the GraphQL schema as returned by the server
// TODO Implement the rest of the schema if needed.
type Schema struct {
	Types []*SchemaType `json:"types"`
}

type SchemaMeta struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Kind        Kind   `json:"kind,omitempty"`
}

// SchemaType defines a specific type within the schema
type SchemaType struct {
	SchemaMeta

	EnumValues    []SchemaEnumValue `json:"enumValues,omitempty"`
	Fields        []SchemaField     `json:"fields,omitempty"`
	InputFields   []SchemaField     `json:"inputFields,omitempty"`
	Interfaces    []SchemaTypeRef   `json:"interfaces,omitempty"`
	PossibleTypes []SchemaTypeRef   `json:"possibleTypes,omitempty"`
}

type SchemaTypeRef struct {
	SchemaMeta

	OfType *SchemaTypeRef `json:"ofType,omitempty"`
}

type SchemaField struct {
	SchemaMeta

	Type         SchemaTypeRef `json:"type"`
	Args         []SchemaField `json:"args,omitempty"`
	DefaultValue interface{}   `json:"defaultValue,omitempty"`
}

type SchemaEnumValue struct {
	SchemaMeta

	IsDeprecated      bool   `json:"isDeprecated"`
	DeprecationReason string `json:"deprecationReason"`
}
