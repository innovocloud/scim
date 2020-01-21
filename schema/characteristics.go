package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// compiled in init in core.go
var attributeNameRegexString = `^[A-Za-z][\w$-]*$`
var attributeNameRegex *regexp.Regexp

func checkAttributeName(name string) error {
	if attributeNameRegex == nil {
		attributeNameRegex = regexp.MustCompile(attributeNameRegexString)
	}
	// starts with a A-Za-z followed by a A-Za-z0-9, a dollar sign, a hyphen or an underscore
	match := attributeNameRegex.MatchString(name)

	if !match {
		return errors.New(fmt.Sprintf("invalid attribute name %q", name)) // @TODO libraries should not panic
	}
	return nil
}

type AttributeMutability int

const (
	AttributeMutabilityReadWrite AttributeMutability = iota
	AttributeMutabilityImmutable
	AttributeMutabilityReadOnly
	AttributeMutabilityWriteOnly
)

func (a AttributeMutability) MarshalJSON() ([]byte, error) {
	switch a {
	case AttributeMutabilityImmutable:
		return json.Marshal("immutable")
	case AttributeMutabilityReadOnly:
		return json.Marshal("readOnly")
	case AttributeMutabilityWriteOnly:
		return json.Marshal("writeOnly")
	default:
		return json.Marshal("readWrite")
	}
}

// AttributeReferenceType is a single keyword indicating the reference type of the SCIM resource that may be referenced.
// This attribute is only applicable for attributes that are of type "reference".
type AttributeReferenceType string

const (
	// AttributeReferenceTypeExternal indicates that the resource is an external resource.
	AttributeReferenceTypeExternal AttributeReferenceType = "external"
	// AttributeReferenceTypeURI indicates that the reference is to a service endpoint or an identifier.
	AttributeReferenceTypeURI   AttributeReferenceType = "uri"
	AttributeReferenceTypeUser  AttributeReferenceType = "User"
	AttributeReferenceTypeGroup AttributeReferenceType = "Group"
)

// AttributeReturned is a single keyword indicating the circumstances under which an attribute and associated values are
// returned in response to a GET request or in response to a PUT, POST, or PATCH request.

type AttributeReturned int

const (
	AttributeReturnedDefault AttributeReturned = iota
	AttributeReturnedAlways
	AttributeReturnedNever
	AttributeReturnedRequest
)

func (a AttributeReturned) MarshalJSON() ([]byte, error) {
	switch a {
	case AttributeReturnedAlways:
		return json.Marshal("always")
	case AttributeReturnedNever:
		return json.Marshal("never")
	case AttributeReturnedRequest:
		return json.Marshal("request")
	default:
		return json.Marshal("default")
	}
}

type DataType int

const (
	DataTypeString DataType = iota
	DataTypeDecimal
	DataTypeInteger
	DataTypeBinary
	DataTypeBoolean
	DataTypeComplex
	DataTypeDateTime
	DataTypeReference
)

func (a DataType) MarshalJSON() ([]byte, error) {
	switch a {
	case DataTypeDecimal:
		return json.Marshal("decimal")
	case DataTypeInteger:
		return json.Marshal("integer")
	case DataTypeBinary:
		return json.Marshal("binary")
	case DataTypeBoolean:
		return json.Marshal("boolean")
	case DataTypeComplex:
		return json.Marshal("complex")
	case DataTypeDateTime:
		return json.Marshal("dateTime")
	case DataTypeReference:
		return json.Marshal("reference")
	default:
		return json.Marshal("string")
	}
}

type AttributeUniqueness int

const (
	AttributeUniquenessNone AttributeUniqueness = iota
	AttributeUniquenessGlobal
	AttributeUniquenessServer
)

func (a AttributeUniqueness) MarshalJSON() ([]byte, error) {
	switch a {
	case AttributeUniquenessGlobal:
		return json.Marshal("global")
	case AttributeUniquenessServer:
		return json.Marshal("server")
	default:
		return json.Marshal("none")
	}
}
