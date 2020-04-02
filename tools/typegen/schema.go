package main

const (
	// https://github.com/graphql/graphql-js/blob/master/src/utilities/getIntrospectionQuery.js#L35
	allTypes = ` query IntrospectionQuery {
      __schema {
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
          ...FullType
        }
        directives {
          name
          description
          locations
          args {
            ...InputValue
          }
        }
      }
    }
    fragment FullType on __Type {
      kind
      name
      description
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }
    fragment InputValue on __InputValue {
      name
      description
      type { ...TypeRef }
      defaultValue
    }
    fragment TypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
	`
)

// Schema contains data about the GraphQL schema as returned by the server
// TODO Implement the rest of the schema if needed.
type Schema struct {
	Types []*SchemaType `json:"types"`
}

// SchemaType defines a specific type within the schema
type SchemaType struct {
	InputFields []SchemaInputValue `json:"inputFields"`
	Kind        string             `json:"kind"`
	Name        string             `json:"name"`
	// Description string             `json:"description"`
	Fields        []SchemaFields    `json:"fields"`
	Interfaces    []SchemaTypeRef   `json:"interfaces"`
	PossibleTypes []SchemaTypeRef   `json:"possibleTypes"`
	EnumValues    []SchemaEnumValue `json:"enumValues"`
}

type SchemaInputValue struct {
	DefaultValue interface{}   `json:"defaultValue"`
	Description  string        `json:"description"`
	Name         string        `json:"name"`
	Type         SchemaTypeRef `json:"type"`
}

type SchemaFields struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Args        []SchemaInputValue `json:"args"`
	Type        SchemaTypeRef      `json:"type"`
}

type SchemaEnumValue struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	IsDeprecated      bool   `json:"isDeprecated"`
	DeprecationReason string `json:"deprecationReason"`
}

type SchemaTypeRef struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	OfType struct {
		Name   string `json:"name"`
		Kind   string `json:"kind"`
		OfType struct {
			Name   string `json:"name"`
			Kind   string `json:"kind"`
			OfType struct {
				Name   string `json:"name"`
				Kind   string `json:"kind"`
				OfType struct {
					Name   string `json:"name"`
					Kind   string `json:"kind"`
					OfType struct {
						Name   string `json:"name"`
						Kind   string `json:"kind"`
						OfType struct {
							Name   string `json:"name"`
							Kind   string `json:"kind"`
							OfType struct {
								Name string `json:"name"`
								Kind string `json:"kind"`
							} `json:"ofType"`
						} `json:"ofType"`
					} `json:"ofType"`
				} `json:"ofType"`
			} `json:"ofType"`
		} `json:"ofType"`
	} `json:"ofType"`
}

type allTypesResponse struct {
	Schema Schema `json:"__schema"`
}
