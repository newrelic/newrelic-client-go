package main

const (
	// https://github.com/graphql/graphql-js/blob/master/src/utilities/getIntrospectionQuery.js#L35
	//   Modified from the following as we only care about the Types
	//   query IntrospectionQuery {
	//     __schema {
	//       directives { name description locations args { ...InputValue } }
	//       mutationType { name }
	//       queryType { name }
	//       subscriptionType { name }
	//       types { ...FullType }
	//     }
	//   }
	allTypes = `
query IntrospectionQuery {
  __schema {
    types { ...FullType }
  }
}
fragment FullType on __Type {
  kind
  name
  description
  fields(includeDeprecated: true) { name description args { ...InputValue } type { ...TypeRef } isDeprecated deprecationReason }
  inputFields { ...InputValue }
  interfaces { ...TypeRef }
  enumValues(includeDeprecated: true) { name description isDeprecated deprecationReason }
  possibleTypes { ...TypeRef }
}
fragment InputValue on __InputValue { name description type { ...TypeRef } defaultValue }
fragment TypeRef on __Type { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name } } } } } } } }
`
)

// allTypesResponse - Util struct for decoding the response
type allTypesResponse struct {
	Schema Schema `json:"__schema"`
}
