package schema

import (
	"strings"

	"github.com/innovocloud/scim/errors"
)

// Schema is a collection of attribute definitions that describe the contents of an entire or partial resource.
type Schema struct {
	Attributes  []CoreAttribute `json:"attributes"`
	Description string          `json:"description,omitempty"`
	ID          string          `json:"id"`
	Name        string          `json:"name,omitempty"`
}

// Validate validates given resource based on the schema.
func (s Schema) Validate(resource interface{}) (map[string]interface{}, errors.ValidationError) {
	core, ok := resource.(map[string]interface{})
	if !ok {
		return nil, errors.ValidationErrorInvalidSyntax
	}

	attributes := make(map[string]interface{})
	for _, attribute := range s.Attributes {
		var hit interface{}
		var found bool
		for k, v := range core {
			if strings.EqualFold(attribute.Name, k) {
				if found {
					return nil, errors.ValidationErrorInvalidSyntax
				}
				found = true
				hit = v
			}
		}

		attr, scimErr := attribute.validate(hit)
		if scimErr != errors.ValidationErrorNil {
			return nil, scimErr
		}
		attributes[attribute.Name] = attr
	}
	return attributes, errors.ValidationErrorNil
}

// ValidatePatchOperationValue validates an individual operation and its related value
func (s Schema) ValidatePatchOperationValue(operation string, operationValue map[string]interface{}) errors.ValidationError {
	for k, v := range operationValue {
		var attr *CoreAttribute
		scimErr := errors.ValidationErrorNil

		for _, attribute := range s.Attributes {
			if strings.EqualFold(attribute.Name, k) {
				attr = &attribute
				break
			}
		}

		// Attribute does not exist in the schema, thus it is an invalid request.
		// Immutable attrs can only be added and Readonly attrs cannot be patched
		if attr == nil || cannotBePatched(operation, *attr) {
			return errors.ValidationErrorInvalidValue
		}

		// "remove" operations simply have to exist
		if operation != "remove" {
			_, scimErr = attr.validate(v)
		}

		if scimErr != errors.ValidationErrorNil {
			return scimErr
		}
	}

	return errors.ValidationErrorNil
}

func cannotBePatched(op string, attr CoreAttribute) bool {
	return isImmutable(op, attr) || isReadOnly(attr)
}

func isImmutable(op string, attr CoreAttribute) bool {
	return attr.Mutability == AttributeMutabilityImmutable && (op == "replace" || op == "remove")
}

func isReadOnly(attr CoreAttribute) bool {
	return attr.Mutability == AttributeMutabilityReadOnly
}

// MarshalJSON converts the schema struct to its corresponding json representation.
// func (s Schema) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(map[string]interface{}{
// 		"id":          s.ID,
// 		"name":        s.Name,
// 		"description": s.Description,
// 		"attributes":  s.getRawAttributes(),
// 	})
// }

// func (s Schema) getRawAttributes() []map[string]interface{} {
// 	attributes := make([]map[string]interface{}, len(s.Attributes))

// 	for i, a := range s.Attributes {
// 		attributes[i] = a.getRawAttributes()
// 	}

// 	return attributes
// }
