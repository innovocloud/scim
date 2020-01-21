package schema

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	datetime "github.com/di-wu/xsd-datetime"
	"github.com/innovocloud/scim/errors"
)

func init() {
	attributeNameRegex = regexp.MustCompile(attributeNameRegexString)
	validStringRegex = regexp.MustCompile(validStringRegexString)
}

// CoreAttribute represents those attributes that sit at the top level of the JSON object together with the common
// attributes (such as the resource "id").
type CoreAttribute struct {
	Name            string                   `json:"name"`
	Type            DataType                 `json:"type"`
	Description     string                   `json:"description,omitempty"`
	MultiValued     bool                     `json:"multiValued"`
	Required        bool                     `json:"required"`
	CaseExact       bool                     `json:"caseExact"`
	CanonicalValues []string                 `json:"canonicalValues,omitempty"`
	Mutability      AttributeMutability      `json:"mutability"`
	Returned        AttributeReturned        `json:"returned"`
	Uniqueness      AttributeUniqueness      `json:"uniqueness"`
	ReferenceTypes  []AttributeReferenceType `json:"referenceTypes,omitempty"`
	SubAttributes   []CoreAttribute          `json:"subAttributes,omitempty"`
}

// SimpleCoreAttribute creates a non-complex attribute based on given parameters.
func SimpleCoreAttribute(params CoreAttribute) CoreAttribute {
	checkAttributeName(params.Name)

	// return CoreAttribute{
	// 	CanonicalValues: params.CanonicalValues,
	// 	CaseExact:       params.CaseExact,
	// 	Description:     params.Description,
	// 	MultiValued:     params.MultiValued,
	// 	Mutability:      params.Mutability,
	// 	Name:            params.Name,
	// 	ReferenceTypes:  params.ReferenceTypes,
	// 	Required:        params.Required,
	// 	Returned:        params.Returned,
	// 	Type:            params.Type,
	// 	Uniqueness:      params.Uniqueness,
	// }
	return CoreAttribute(params)
}

// ComplexCoreAttribute creates a complex attribute based on given parameters.
func ComplexCoreAttribute(params CoreAttribute) CoreAttribute {
	checkAttributeName(params.Name)

	names := map[string]int{}
	var sa []CoreAttribute
	for i, a := range params.SubAttributes {
		name := strings.ToLower(a.Name)
		if j, ok := names[name]; ok {
			panic(fmt.Errorf("duplicate name %q for sub-attributes %d and %d", name, i, j))
		}
		names[name] = i

		sa = append(sa, CoreAttribute(a))
	}

	return CoreAttribute{
		Name:          params.Name,
		Description:   params.Description,
		MultiValued:   params.MultiValued,
		Mutability:    params.Mutability,
		Required:      params.Required,
		Returned:      params.Returned,
		SubAttributes: sa,
		Type:          DataTypeComplex,
		Uniqueness:    params.Uniqueness,
	}
}

// type ComplexParams CoreAttribute

func (a CoreAttribute) validate(attribute interface{}) (interface{}, errors.ValidationError) {
	// return false if the attribute is not present but required.
	if attribute == nil {
		if !a.Required {
			return nil, errors.ValidationErrorNil
		}
		return nil, errors.ValidationErrorInvalidValue
	}

	if a.MultiValued {
		// return false if the multivalued attribute is not a slice.
		arr, ok := attribute.([]interface{})
		if !ok {
			return nil, errors.ValidationErrorInvalidSyntax
		}

		// return false if the multivalued attribute is empty.
		if a.Required && len(arr) == 0 {
			return nil, errors.ValidationErrorInvalidValue
		}

		attributes := make([]interface{}, 0)
		for _, ele := range arr {
			attr, scimErr := a.validateSingular(ele)
			if scimErr != errors.ValidationErrorNil {
				return nil, scimErr
			}
			attributes = append(attributes, attr)
		}
		return attributes, errors.ValidationErrorNil
	}

	return a.validateSingular(attribute)
}

// compiled in init at the top of the file
var validStringRegexString = `^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`
var validStringRegex *regexp.Regexp

func (a CoreAttribute) validateSingular(attribute interface{}) (interface{}, errors.ValidationError) {
	switch a.Type {
	case DataTypeBinary:
		bin, ok := attribute.(string)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}

		match := validStringRegex.MatchString(bin)
		if !match {
			return nil, errors.ValidationErrorInvalidValue
		}

		return bin, errors.ValidationErrorNil
	case DataTypeBoolean:
		b, ok := attribute.(bool)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}
		return b, errors.ValidationErrorNil
	case DataTypeComplex:
		complex, ok := attribute.(map[string]interface{})
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}

		attributes := make(map[string]interface{})
		for _, sub := range a.SubAttributes {
			var hit interface{}
			var found bool
			for k, v := range complex {
				if strings.EqualFold(sub.Name, k) {
					if found {
						return nil, errors.ValidationErrorInvalidSyntax
					}
					found = true
					hit = v
				}
			}

			attr, scimErr := sub.validate(hit)
			if scimErr != errors.ValidationErrorNil {
				return nil, scimErr
			}
			attributes[sub.Name] = attr
		}
		return attributes, errors.ValidationErrorNil
	case DataTypeDateTime:
		date, ok := attribute.(string)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}
		_, err := datetime.Parse(date)
		if err != nil {
			return nil, errors.ValidationErrorInvalidValue
		}
		return date, errors.ValidationErrorNil
	case DataTypeDecimal:
		if reflect.TypeOf(attribute).Kind() != reflect.Float64 {
			return nil, errors.ValidationErrorInvalidValue
		}
		return attribute.(float64), errors.ValidationErrorNil
	case DataTypeInteger:
		if reflect.TypeOf(attribute).Kind() != reflect.Int {
			return nil, errors.ValidationErrorInvalidValue
		}
		return attribute.(int), errors.ValidationErrorNil
	case DataTypeString, DataTypeReference:
		s, ok := attribute.(string)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}
		return s, errors.ValidationErrorNil
	default:
		return nil, errors.ValidationErrorInvalidSyntax
	}
}
