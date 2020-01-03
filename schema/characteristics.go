package schema

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func checkAttributeName(name string) {
	// starts with a A-Za-z followed by a A-Za-z0-9, a dollar sign, a hyphen or an underscore
	match, err := regexp.MatchString(`^[A-Za-z][\w$-]*$`, name)
	if err != nil {
		panic(err) // @TODO libraries should not panic
	}

	if !match {
		panic(fmt.Sprintf("invalid attribute name %q", name)) // @TODO libraries should not panic
	}
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
	AttributeReferenceTypeURI AttributeReferenceType = "uri"
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

type AttributeDataType int

const (
	AttributeDataTypeDecimal AttributeDataType = iota
	AttributeDataTypeInteger

	AttributeDataTypeBinary
	AttributeDataTypeBoolean
	AttributeDataTypeComplex
	AttributeDataTypeDateTime
	AttributeDataTypeReference
	AttributeDataTypeString
)

func (a AttributeDataType) MarshalJSON() ([]byte, error) {
	switch a {
	case AttributeDataTypeDecimal:
		return json.Marshal("decimal")
	case AttributeDataTypeInteger:
		return json.Marshal("integer")
	case AttributeDataTypeBinary:
		return json.Marshal("binary")
	case AttributeDataTypeBoolean:
		return json.Marshal("boolean")
	case AttributeDataTypeComplex:
		return json.Marshal("complex")
	case AttributeDataTypeDateTime:
		return json.Marshal("dateTime")
	case AttributeDataTypeReference:
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
