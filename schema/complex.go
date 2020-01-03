package schema

// ComplexParams are the parameters used to create a complex attribute.
type ComplexParams struct {
	Description   string `json:",omitempty"`
	MultiValued   bool
	Mutability    AttributeMutability
	Name          string
	Required      bool
	Returned      AttributeReturned
	SubAttributes []SimpleParams
	Uniqueness    AttributeUniqueness
}
