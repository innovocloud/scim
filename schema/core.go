package schema

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	datetime "github.com/di-wu/xsd-datetime"
	"github.com/innovocloud/scim/errors"
)

// SimpleCoreAttribute creates a non-complex attribute based on given parameters.
func SimpleCoreAttribute(params SimpleParams) CoreAttribute {
	checkAttributeName(params.name)

	return CoreAttribute{
		CanonicalValues: params.canonicalValues,
		CaseExact:       params.caseExact,
		Description:     params.Description,
		MultiValued:     params.multiValued,
		Mutability:      params.mutability,
		Name:            params.name,
		ReferenceTypes:  params.referenceTypes,
		Required:        params.required,
		Returned:        params.returned,
		Typ:             params.typ,
		Uniqueness:      params.uniqueness,
	}
}

// ComplexCoreAttribute creates a complex attribute based on given parameters.
func ComplexCoreAttribute(params ComplexParams) CoreAttribute {
	checkAttributeName(params.Name)

	names := map[string]int{}
	var sa []CoreAttribute
	for i, a := range params.SubAttributes {
		name := strings.ToLower(a.name)
		if j, ok := names[name]; ok {
			panic(fmt.Errorf("duplicate name %q for sub-attributes %d and %d", name, i, j))
		}
		names[name] = i

		sa = append(sa, CoreAttribute{
			CanonicalValues: a.canonicalValues,
			CaseExact:       a.caseExact,
			Description:     a.Description,
			MultiValued:     a.multiValued,
			Mutability:      a.mutability,
			Name:            a.name,
			ReferenceTypes:  a.referenceTypes,
			Required:        a.required,
			Returned:        a.returned,
			Typ:             a.typ,
			Uniqueness:      a.uniqueness,
		})
	}

	return CoreAttribute{
		Name:          params.Name,
		Description:   params.Description,
		MultiValued:   params.MultiValued,
		Mutability:    params.Mutability,
		Required:      params.Required,
		Returned:      params.Returned,
		subAttributes: sa,
		Typ:           AttributeDataTypeComplex,
		Uniqueness:    params.Uniqueness,
	}
}

// CoreAttribute represents those attributes that sit at the top level of the JSON object together with the common
// attributes (such as the resource "id").
type CoreAttribute struct {
	Name            string
	Typ             AttributeDataType `json:"type"`
	MultiValued     bool
	Description     string `json:",omitempty"`
	Required        bool
	CaseExact       bool
	CanonicalValues []string
	Mutability      AttributeMutability
	Returned        AttributeReturned
	Uniqueness      AttributeUniqueness
	ReferenceTypes  []AttributeReferenceType
	subAttributes   []CoreAttribute
}

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

func (a CoreAttribute) validateSingular(attribute interface{}) (interface{}, errors.ValidationError) {
	switch a.Typ {
	case AttributeDataTypeBinary:
		bin, ok := attribute.(string)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}

		match, err := regexp.MatchString(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`, bin)
		if err != nil {
			panic(err)
		}

		if !match {
			return nil, errors.ValidationErrorInvalidValue
		}

		return bin, errors.ValidationErrorNil
	case AttributeDataTypeBoolean:
		b, ok := attribute.(bool)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}
		return b, errors.ValidationErrorNil
	case AttributeDataTypeComplex:
		complex, ok := attribute.(map[string]interface{})
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}

		attributes := make(map[string]interface{})
		for _, sub := range a.subAttributes {
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
	case AttributeDataTypeDateTime:
		date, ok := attribute.(string)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}
		_, err := datetime.Parse(date)
		if err != nil {
			return nil, errors.ValidationErrorInvalidValue
		}
		return date, errors.ValidationErrorNil
	case AttributeDataTypeDecimal:
		if reflect.TypeOf(attribute).Kind() != reflect.Float64 {
			return nil, errors.ValidationErrorInvalidValue
		}
		return attribute.(float64), errors.ValidationErrorNil
	case AttributeDataTypeInteger:
		if reflect.TypeOf(attribute).Kind() != reflect.Int {
			return nil, errors.ValidationErrorInvalidValue
		}
		return attribute.(int), errors.ValidationErrorNil
	case AttributeDataTypeString, AttributeDataTypeReference:
		s, ok := attribute.(string)
		if !ok {
			return nil, errors.ValidationErrorInvalidValue
		}
		return s, errors.ValidationErrorNil
	default:
		return nil, errors.ValidationErrorInvalidSyntax
	}
}

func (a *CoreAttribute) getRawAttributes() map[string]interface{} {
	rawSubAttributes := make([]map[string]interface{}, len(a.subAttributes))

	for i, subAttr := range a.subAttributes {
		rawSubAttributes[i] = subAttr.getRawAttributes()
	}

	return map[string]interface{}{
		"canonicalValues": a.CanonicalValues,
		"caseExact":       a.CaseExact,
		"description":     a.Description,
		"multiValued":     a.MultiValued,
		"mutability":      a.Mutability,
		"name":            a.Name,
		"referenceTypes":  a.ReferenceTypes,
		"required":        a.Required,
		"returned":        a.Returned,
		"subAttributes":   rawSubAttributes,
		"type":            a.Typ,
		"uniqueness":      a.Uniqueness,
	}
}
