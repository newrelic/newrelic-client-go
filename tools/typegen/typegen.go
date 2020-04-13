package main

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (r *SchemaTypeRef) IsList() bool {
	kinds := r.GetKinds()

	if len(kinds) > 0 && kinds[0] == KindList {
		return true
	}

	return false
}

// GetKind returns an array or the type kind
func (r *SchemaTypeRef) GetKinds() []Kind {
	tree := []Kind{}

	if r.Kind != "" && r.Kind != KindNonNull {
		tree = append(tree, r.Kind)
	}

	// Recursion FTW
	if r.OfType != nil {
		tree = append(tree, r.OfType.GetKinds()...)
	}

	return tree
}

// GetName returns a recusive lookup of the type name
func (m *SchemaMeta) GetName() string {
	var fieldName string

	switch strings.ToLower(m.Name) {
	case "ids":
		// special case to avoid the struct field Ids, and prefer IDs instead
		fieldName = "IDs"
	case "id":
		fieldName = "ID"
	case "accountid":
		fieldName = "AccountID"
	default:
		fieldName = strings.Title(m.Name)
	}

	return fieldName
}

func (m *SchemaMeta) GetTags() string {
	if m == nil {
		return ""
	}

	jsonTag := "`json:\"" + m.Name

	// Overrides
	if strings.EqualFold(m.Name, "id") {
		jsonTag += ",string"
	}

	return jsonTag + "\"`"
}

// GetName returns a recusive lookup of the type name
func (r *SchemaTypeRef) GetTypeName() string {
	if r.Name != "" {
		return r.Name
	}

	// Recursion FTW
	if r.OfType != nil {
		return r.OfType.GetTypeName()
	}

	log.Errorf("failed to get name for %#v", *r)
	return "UNKNOWN"
}

// FieldType resolves the given SchemaInputField into a field name to use on a go struct.
//  type, recurse, error
func (r *SchemaTypeRef) GetType() (string, bool, error) {
	switch n := r.GetTypeName(); n {
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
		return "", true, fmt.Errorf("empty field name: %+v", r)
	default:
		return n, true, nil
	}
}

// GetDescription looks for anything in the description before \n\n---\n
// and filters off anything after that (internal messaging that is not useful here)
func (m *SchemaMeta) GetDescription() string {
	var ret string

	if strings.TrimSpace(m.Description) == "" {
		return ""
	}

	r := regexp.MustCompile(`(?s)(.*)\n---\n`)
	desc := r.FindStringSubmatch(m.Description)

	log.Debugf("Description: %#v", desc)

	if len(desc) > 1 {
		ret = desc[1]
	} else {
		ret = m.Description
	}

	return "\t /* " + m.GetName() + " - " + strings.TrimSpace(ret) + " */\n"
}

// Global type list lookup function
func (s *Schema) LookupTypeByName(typeName string) (*SchemaType, error) {
	log.Debugf("looking for typeName: %s", typeName)

	for _, t := range s.Types {
		if t.Name == typeName {
			return t, nil
		}
	}

	return nil, fmt.Errorf("type by name %s not found", typeName)
}

// Definition generates the Golang definition of the type
func (s *Schema) Definition(typeInfo TypeConfig) (string, error) {
	t, err := s.LookupTypeByName(typeInfo.Name)
	if err != nil {
		return "", err
	}

	// Start with the type description
	output := t.GetDescription()

	switch t.Kind {
	case KindInputObject, KindObject:
		output += "type " + t.Name + " struct {\n"

		// Fill in the struct fields for an input type
		for _, f := range t.InputFields {
			output += s.lineForField(f)
		}

		for _, f := range t.Fields {
			output += s.lineForField(f)
		}

		output += "}\n"
	case KindENUM:
		output += "type " + t.Name + " string\n\n"
		output += "const (\n"

		for _, v := range t.EnumValues {
			output += v.GetDescription()
			output += "\t" + v.Name + " " + t.Name + " = \"" + v.Name + "\"\n"
		}

		output += ")\n"
	case KindScalar:
		// Default to string for scalars, but warn this is might not be what they want.
		createAs := "string"
		if typeInfo.CreateAs != "" {
			createAs = typeInfo.CreateAs
		} else {
			log.Warnf("creating scalar %s as string", t.Name)
		}

		output += "type " + t.Name + " " + createAs + "\n"
	case KindInterface:
		createAs := "interface{}"
		if typeInfo.CreateAs != "" {
			createAs = typeInfo.CreateAs
		}

		output += "type " + t.Name + " " + createAs + "\n"

	default:
		log.Warnf("unhandled object Kind: %s\n", t.Kind)
	}

	return output + "\n", nil
}

func (s *Schema) lineForField(f SchemaField) string {
	output := f.GetDescription()

	log.Infof("handling kind %s: %+v", f.Type.Kind, f.Type)
	fieldType, recurse, err := f.Type.GetType()
	if err != nil {
		// If we have an error, then we don't know how to handle the type to
		// determine the field name.
		log.Errorf("error resolving first non-empty name from field: %#v: %s", f.Type, err.Error())
	}

	if recurse {
		log.Debugf("recurse search for %s: %+v", fieldType, f.Type)

		// The name of the nested sub-type.  We take the first value here as the root name for the nested type.
		subTName := f.Type.GetTypeName()
		log.Tracef("subTName %s", subTName)

		err := s.TypeGen(TypeConfig{Name: subTName})
		if err != nil {
			log.Errorf("ERROR while resolving sub type %s: %s\n", subTName, err)
		}

		fieldType = subTName
	}

	fieldTypePrefix := ""

	if f.Type.IsList() {
		fieldTypePrefix = "[]"
	}

	fieldTags := f.GetTags()

	output += "\t" + f.GetName() + " " + fieldTypePrefix + fieldType + " " + fieldTags + "\n"

	return output
}

// TypeGen is the mother type generator.
func (s *Schema) TypeGen(typeInfo TypeConfig) error {
	log.Infof("starting on: %+v", typeInfo)

	// Only add the new types
	if _, ok := types[typeInfo.Name]; !ok {
		output, err := s.Definition(typeInfo)
		if err != nil {
			return err
		}

		types[typeInfo.Name] = output
	}

	return nil
}
