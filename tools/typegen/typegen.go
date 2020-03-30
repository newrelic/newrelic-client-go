package main

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	types = make(map[string]string)
)

func ResolveSchemaTypes(schema Schema, typeNames []string) (map[string]string, error) {
	typeKeeper := make(map[string]string)

	log.SetLevel(log.DebugLevel)

	for _, typeName := range typeNames {
		typeGenResult, err := TypeGen(schema, typeName)
		if err != nil {
			log.Errorf("error while generating type %s: %s", typeName, err)
		}

		for k, v := range typeGenResult {
			typeKeeper[k] = v
		}
	}

	// collect the types stored from recursion also.
	for k, v := range types {
		typeKeeper[k] = v
	}

	return typeKeeper, nil
}

func handleEnumType(schema Schema, t SchemaType) map[string]string {
	typeKeeper := make(map[string]string)

	// output collects each line of a struct type
	output := []string{}

	// Add a comment for golint to ignore
	output = append(output, "")

	output = append(output, "// nolint:golint")
	output = append(output, fmt.Sprintf("type %s string ", t.Name))
	output = append(output, "")

	output = append(output, "const (")
	for _, v := range t.EnumValues {

		if v.Description != "" {
			output = append(output, fmt.Sprintf("\t /* %s */", parseDescription(v.Description)))
		}

		output = append(output, fmt.Sprintf("\t%s %s = \"%s\" // nolint:golint", v.Name, t.Name, v.Name))
	}

	output = append(output, ")")
	output = append(output, "")

	typeKeeper[t.Name] = strings.Join(output, "\n")

	return typeKeeper
}

func handleScalarType(schema Schema, t SchemaType) map[string]string {
	typeKeeper := make(map[string]string)

	// output collects each line of a struct type
	output := []string{}

	// Add a comment for golint to ignore
	output = append(output, "")

	output = append(output, "// nolint:golint")
	output = append(output, fmt.Sprintf("type %s string ", t.Name))
	output = append(output, "")

	typeKeeper[t.Name] = strings.Join(output, "\n")

	return typeKeeper
}

func kindTree(t SchemaTypeRef) []string {
	tree := []string{}

	if t.Kind != "" {
		tree = append(tree, t.Kind)
	}

	if t.OfType.Kind != "" {
		tree = append(tree, t.OfType.Kind)
	}

	if t.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.OfType.Kind)
	}

	return tree
}

func nameTree(t SchemaTypeRef) []string {
	tree := []string{}

	if t.Name != "" {
		tree = append(tree, t.Name)
	}

	if t.OfType.Name != "" {
		tree = append(tree, t.OfType.Name)
	}

	if t.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.OfType.Name)
	}

	return tree
}

func removeNonNullValues(tree []string) []string {
	a := []string{}

	for _, x := range tree {
		if x != "NON_NULL" {
			a = append(a, x)
		}
	}

	return a
}

// fieldTypeFromTypeRef resolves the given SchemaInputValue into a field name to use on a go struct.
func fieldTypeFromTypeRef(t SchemaTypeRef) (string, bool, error) {

	switch n := nameTree(t)[0]; n {
	case "String":
		return "string", false, nil
	case "Int":
		return "int", false, nil
	case "Boolean":
		return "bool", false, nil
	case "Float":
		return "float64", false, nil
	case "ID":
		// ID is a nested object, but behaves like an integer.  This may be true of other SCALAR types as well, so logic here could potentially be moved.
		return "int", false, nil
	case "":
		return "", true, fmt.Errorf("empty field name: %+v", t)
	default:
		return n, true, nil
	}
}

// handleObjectType will operate on a SchemaType who's Kind is OBJECT or INPUT_OBJECT.
func handleObjectType(schema Schema, t SchemaType) map[string]string {
	typeKeeper := make(map[string]string)

	// output collects each line of a struct type
	output := []string{}

	// Add a comment for golint to ignore
	output = append(output, "// nolint:golint")

	output = append(output, fmt.Sprintf("type %s struct {", t.Name))

	// Fill in the struct fields for an input type
	for _, f := range t.InputFields {
		log.Debugf("Input Field: %+v", f.Name)
		output = append(output, "")
		output = append(output, lineForField(schema, f.Name, f.Description, f.Type)...)
	}

	for _, f := range t.Fields {
		log.Debugf("Field: %+v", f.Name)
		output = append(output, "")
		output = append(output, lineForField(schema, f.Name, f.Description, f.Type)...)
	}

	// Close the struct
	output = append(output, "}\n")
	typeKeeper[t.Name] = strings.Join(output, "\n")

	return typeKeeper
}

func lineForField(schema Schema, name string, description string, typeRef SchemaTypeRef) []string {
	var output []string
	var fieldName string

	log.Infof("handling kind %s: %+v", typeRef.Kind, typeRef)
	fieldType, recurse, err := fieldTypeFromTypeRef(typeRef)
	if err != nil {
		// If we have an error, then we don't know how to handle the type to
		// determine the field name.
		log.Errorf("error resolving first non-empty name from field: %s: %s", typeRef, err)
	}

	if recurse {
		log.Debugf("recurse search for %s: %+v", fieldType, typeRef)

		// The name of the nested sub-type.  We take the first value here as the root name for the nested type.
		subTName := nameTree(typeRef)[0]

		log.Debugf("subTName %+v", subTName)

		subT, err := typeByName(schema, subTName)
		if err != nil {
			log.Warnf("non_null: unhandled type: %+v\n", name)
			// break
		}

		// Determnine if we need to resolve the sub type, or if it already
		// exists in the map.
		if _, ok := types[subT.Name]; !ok {
			result, err := TypeGen(schema, subT.Name)
			if err != nil {
				log.Errorf("ERROR while resolving sub type %s: %s\n", subT.Name, err)
			}

			log.Debugf("resolved type result:\n%+v\n", result)

			for k, v := range result {
				if _, ok := types[k]; !ok {
					types[k] = v
				}
			}
		}

		fieldType = subT.Name
	}

	if name == "ids" {
		// special case to avoid the struct field Ids, and prefer IDs instead
		fieldName = "IDs"
	} else if name == "id" {
		fieldName = "ID"
	} else if name == "accountId" {
		fieldName = "AccountID"
	} else {
		fieldName = strings.Title(name)
	}

	fieldTypePrefix := ""

	if removeNonNullValues(kindTree(typeRef))[0] == "LIST" {
		fieldTypePrefix = "[]"
	}

	// Include some documentation
	if description != "" {
		output = append(output, "\t /* "+parseDescription(description)+" */")
	}

	var fieldTags string
	if name == "id" {
		fieldTags = fmt.Sprintf("`json:\"%s,string\"`", name)
	} else {
		fieldTags = fmt.Sprintf("`json:\"%s\"`", name)
	}

	output = append(output, fmt.Sprintf("\t %s %s%s %s", fieldName, fieldTypePrefix, fieldType, fieldTags))

	return output
}

// parseDescription looks for anything in the description before \n\n---\n
// and filters off anything after that (internal messaging that is not useful here)
func parseDescription(description string) string {
	//r := regexp.MustCompile(`(.*)\r\n---\r\n.*`)
	r := regexp.MustCompile(`(?s)(.*)\n---\n`)
	desc := r.FindStringSubmatch(description)

	log.Debugf("Description: %#v", desc)

	if len(desc) > 1 {
		if strings.Count(desc[1], "\n") < 2 {
			return strings.Trim(desc[1], "\n")
		}
		return desc[1]
	}

	return description
}

// TypeGen is the mother type generator.
func TypeGen(schema Schema, typeName string) (map[string]string, error) {

	// The total known types.  Keyed by the typeName, and valued as the string
	// output that one would write to a file where Go structs are kept.
	typeKeeper := make(map[string]string)

	t, err := typeByName(schema, typeName)
	if err != nil {
		log.Error(err)
	}

	log.Infof("starting on %s: %+v", typeName, t.Kind)

	// To store the results from the single
	results := make(map[string]string)

	if t.Kind == "INPUT_OBJECT" || t.Kind == "OBJECT" {
		results = handleObjectType(schema, *t)
	} else if t.Kind == "ENUM" {
		results = handleEnumType(schema, *t)
	} else if t.Kind == "SCALAR" {
		results = handleScalarType(schema, *t)
	} else {
		log.Warnf("WARN: unhandled object Kind: %s\n", t.Kind)
	}

	for k, v := range results {
		typeKeeper[k] = v
	}

	// return strings.Join(output, "\n"), nil
	return typeKeeper, nil
}

func typeByName(schema Schema, typeName string) (*SchemaType, error) {
	log.Debugf("looking for typeName: %s", typeName)

	for _, t := range schema.Types {
		if t.Name == typeName {
			return t, nil
		}
	}

	return nil, fmt.Errorf("type by name %s not found", typeName)
}
