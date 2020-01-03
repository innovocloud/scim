package schema

// SimpleParams are the parameters used to create a simple attribute.
type SimpleParams struct {
	canonicalValues []string
	caseExact       bool
	Description     string `json:",omitempty"`
	multiValued     bool
	mutability      AttributeMutability
	name            string
	referenceTypes  []AttributeReferenceType
	required        bool
	returned        AttributeReturned
	typ             AttributeDataType
	uniqueness      AttributeUniqueness
}

// SimpleBinaryParams converts given binary parameters to their corresponding simple parameters.
func SimpleBinaryParams(params BinaryParams) SimpleParams {
	return SimpleParams{
		caseExact:   true,
		Description: params.Description,
		multiValued: params.MultiValued,
		mutability:  params.Mutability,
		name:        params.Name,
		required:    params.Required,
		returned:    params.Returned,
		typ:         AttributeDataTypeBinary,
		uniqueness:  AttributeUniquenessNone,
	}
}

// BinaryParams are the parameters used to create a simple attribute with a data type of "binary".
// The attribute value MUST be base64 encoded. In JSON representation, the encoded values are represented as a JSON string.
// A binary is case exact and has no uniqueness.
type BinaryParams struct {
	Description string `json:",omitempty"`
	MultiValued bool
	Mutability  AttributeMutability
	Name        string
	Required    bool
	Returned    AttributeReturned
}

// SimpleBooleanParams converts given boolean parameters to their corresponding simple parameters.
func SimpleBooleanParams(params BooleanParams) SimpleParams {
	return SimpleParams{
		caseExact:   false,
		Description: params.Description,
		multiValued: params.MultiValued,
		mutability:  params.Mutability,
		name:        params.Name,
		required:    params.Required,
		returned:    params.Returned,
		typ:         AttributeDataTypeBoolean,
		uniqueness:  AttributeUniquenessNone,
	}
}

// BooleanParams are the parameters used to create a simple attribute with a data type of "boolean".
// The literal "true" or "false". A boolean has no case sensitivity or uniqueness.
type BooleanParams struct {
	Description string `json:",omitempty"`
	MultiValued bool
	Mutability  AttributeMutability
	Name        string
	Required    bool
	Returned    AttributeReturned
}

// SimpleDateTimeParams converts given date time parameters to their corresponding simple parameters.
func SimpleDateTimeParams(params DateTimeParams) SimpleParams {
	return SimpleParams{
		caseExact:   false,
		Description: params.Description,
		multiValued: params.MultiValued,
		mutability:  params.Mutability,
		name:        params.Name,
		required:    params.Required,
		returned:    params.Returned,
		typ:         AttributeDataTypeDateTime,
		uniqueness:  AttributeUniquenessNone,
	}
}

// DateTimeParams are the parameters used to create a simple attribute with a data type of "dateTime".
// A DateTime value (e.g., 2008-01-23T04:56:22Z). A date time format has no case sensitivity or uniqueness.
type DateTimeParams struct {
	Description string `json:",omitempty"`
	MultiValued bool
	Mutability  AttributeMutability
	Name        string
	Required    bool
	Returned    AttributeReturned
}

// SimpleNumberParams converts given number parameters to their corresponding simple parameters.
func SimpleNumberParams(params NumberParams) SimpleParams {
	return SimpleParams{
		caseExact:   false,
		Description: params.Description,
		multiValued: params.MultiValued,
		mutability:  params.Mutability,
		name:        params.Name,
		required:    params.Required,
		returned:    params.Returned,
		typ:         params.Type,
		uniqueness:  params.Uniqueness,
	}
}

// NumberParams are the parameters used to create a simple attribute with a data type of "decimal" or "integer".
// A number has no case sensitivity.
type NumberParams struct {
	Description string `json:",omitempty"`
	MultiValued bool
	Mutability  AttributeMutability
	Name        string
	Required    bool
	Returned    AttributeReturned
	Type        AttributeDataType
	Uniqueness  AttributeUniqueness
}

// SimpleReferenceParams converts given reference parameters to their corresponding simple parameters.
func SimpleReferenceParams(params ReferenceParams) SimpleParams {
	return SimpleParams{
		caseExact:      true,
		Description:    params.Description,
		multiValued:    params.MultiValued,
		mutability:     params.Mutability,
		name:           params.Name,
		referenceTypes: params.ReferenceTypes,
		required:       params.Required,
		returned:       params.Returned,
		typ:            AttributeDataTypeReference,
		uniqueness:     params.Uniqueness,
	}
}

// ReferenceParams are the parameters used to create a simple attribute with a data type of "reference".
// A reference is case exact. A reference has a "referenceTypes" attribute that indicates what types of resources may
// be linked.
type ReferenceParams struct {
	Description    string `json:",omitempty"`
	MultiValued    bool
	Mutability     AttributeMutability
	Name           string
	ReferenceTypes []AttributeReferenceType
	Required       bool
	Returned       AttributeReturned
	Uniqueness     AttributeUniqueness
}

// SimpleStringParams converts given string parameters to their corresponding simple parameters.
func SimpleStringParams(params StringParams) SimpleParams {
	return SimpleParams{
		canonicalValues: params.CanonicalValues,
		caseExact:       params.CaseExact,
		Description:     params.Description,
		multiValued:     params.MultiValued,
		mutability:      params.Mutability,
		name:            params.Name,
		required:        params.Required,
		returned:        params.Returned,
		typ:             AttributeDataTypeString,
		uniqueness:      params.Uniqueness,
	}
}

// StringParams are the parameters used to create a simple attribute with a data type of "string".
// A string is a sequence of zero or more Unicode characters encoded using UTF-8.
type StringParams struct {
	CanonicalValues []string
	CaseExact       bool
	Description     string `json:",omitempty"`
	MultiValued     bool
	Mutability      AttributeMutability
	Name            string
	Required        bool
	Returned        AttributeReturned
	Uniqueness      AttributeUniqueness
}
