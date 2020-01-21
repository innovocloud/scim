package schema

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/innovocloud/scim/errors"
)

func TestInvalidAttributeName(t *testing.T) {
	t.Skip() // @TODO
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	_ = Schema{
		ID:          "urn:ietf:params:scim:schemas:core:2.0:User",
		Name:        "User",
		Description: "User Account",
		Attributes: []CoreAttribute{
			SimpleCoreAttribute(CoreAttribute{Name: "_Invalid"}),
		},
	}
}

var testSchema = Schema{
	ID:          "empty",
	Name:        "empty",
	Description: "",
	Attributes: []CoreAttribute{
		CoreAttribute{
			Name:     "required",
			Required: true,
		},
		CoreAttribute{
			MultiValued: true,
			Name:        "booleans",
			Required:    true,
			Type:        DataTypeBoolean,
		},
		ComplexCoreAttribute(CoreAttribute{
			MultiValued: true,
			Name:        "complex",
			SubAttributes: []CoreAttribute{
				CoreAttribute{Name: "sub"},
			},
		}),

		CoreAttribute{
			Name:      "binary",
			Type:      DataTypeBinary,
			CaseExact: true,
		},
		CoreAttribute{
			Name: "dateTime",
			Type: DataTypeDateTime,
		},
		CoreAttribute{
			Name:      "reference",
			Type:      DataTypeReference,
			CaseExact: true,
		},
		CoreAttribute{
			Name: "integer",
			Type: DataTypeInteger,
		},
		CoreAttribute{
			Name: "decimal",
			Type: DataTypeDecimal,
		},
	},
}

func TestResourceInvalid(t *testing.T) {
	var resource interface{}
	if _, scimErr := testSchema.Validate(resource); scimErr == errors.ValidationErrorNil {
		t.Error("invalid resource expected")
	}
}

func TestValidationInvalid(t *testing.T) {
	for _, test := range []map[string]interface{}{
		{ // missing required field
			"field": "present",
			"booleans": []interface{}{
				true,
			},
		},
		{ // missing required multivalued field
			"required": "present",
			"booleans": []interface{}{},
		},
		{ // wrong type element of slice
			"required": "present",
			"booleans": []interface{}{
				"present",
			},
		},
		{ // duplicate names
			"required": "present",
			"Required": "present",
			"booleans": []interface{}{
				true,
			},
		},
		{ // wrong string type
			"required": true,
			"booleans": []interface{}{
				true,
			},
		},
		{ // wrong complex type
			"required": "present",
			"complex":  "present",
			"booleans": []interface{}{
				true,
			},
		},
		{ // wrong complex element type
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				"present",
			},
		},
		{ // duplicate complex element names
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				map[string]interface{}{
					"sub": "present",
					"Sub": "present",
				},
			},
		},
		{ // wrong type complex element
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				map[string]interface{}{
					"sub": true,
				},
			},
		},
		{ // invalid type binary
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"binary": true,
		},
		{ // invalid type dateTime
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"dateTime": "04:56:22Z2008-01-23T",
		},
		{ // invalid type integer
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"integer": 1.1,
		},
		{ // invalid type decimal
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"decimal": "1.1",
		},
	} {
		if _, scimErr := testSchema.Validate(test); scimErr == errors.ValidationErrorNil {
			t.Errorf("invalid resource expected")
		}
	}
}

func TestValidValidation(t *testing.T) {
	for _, test := range []map[string]interface{}{
		{
			"required": "present",
			"booleans": []interface{}{
				true,
			},
			"complex": []interface{}{
				map[string]interface{}{
					"sub": "present",
				},
			},
			"binary":   "ZXhhbXBsZQ==",
			"dateTime": "2008-01-23T04:56:22Z",
			"integer":  11,
			"decimal":  -2.1e5,
		},
	} {
		if _, scimErr := testSchema.Validate(test); scimErr != errors.ValidationErrorNil {
			t.Errorf("valid resource expected")
		}
	}
}

func TestJSONMarshalling(t *testing.T) {
	expectedJSON, err := ioutil.ReadFile("./fixtures/schema_test.json")

	if err != nil {
		t.Errorf("Failed to required test fixture")
		return
	}

	actualJSON, err := json.Marshal(testSchema)

	if err != nil {
		t.Errorf("Failed to marshal schema into JSON")
		return
	}

	normalizedActual, err := normalizeJSON(actualJSON)
	normalizedExpected, expectedErr := normalizeJSON(expectedJSON)

	if err != nil || expectedErr != nil {
		t.Errorf("Failed to normalize test JSON")
		return
	}

	if normalizedActual != normalizedExpected {
		t.Errorf("Schema output by MarshalJSON did not match the expected output. Want:\n%s\nGot:\n%s\n", normalizedExpected, normalizedActual)
	}
}

func normalizeJSON(rawJSON []byte) (string, error) {
	dataMap := map[string]interface{}{}

	err := json.Unmarshal(rawJSON, &dataMap)
	if err != nil {
		return "", err
	}

	ret, err := json.Marshal(dataMap)

	return string(ret), err
}
