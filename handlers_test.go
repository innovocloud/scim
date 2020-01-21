package scim

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/innovocloud/scim/schema"
)

func newTestServer() Server {
	userSchema := schema.Schema{
		ID:          "urn:ietf:params:scim:schemas:core:2.0:User",
		Name:        "User",
		Description: "User Account",
		Attributes: []schema.CoreAttribute{
			{
				Name:       "userName",
				Required:   true,
				Uniqueness: schema.AttributeUniquenessServer,
			},
			{
				Name:     "active",
				Required: false,
				Type:     schema.DataTypeBoolean,
			},
			{
				Name:       "readonlyThing",
				Required:   false,
				Mutability: schema.AttributeMutabilityReadOnly,
			},
			{
				Name:       "immutableThing",
				Required:   false,
				Mutability: schema.AttributeMutabilityImmutable,
			},
			schema.ComplexCoreAttribute(schema.CoreAttribute{
				Name:     "Name",
				Required: false,
				SubAttributes: []schema.CoreAttribute{
					{
						Name: "familyName",
					},
					{
						Name: "givenName",
					},
				},
			}),
			{
				Name: "displayName",
			},
			schema.ComplexCoreAttribute(schema.CoreAttribute{
				Name:        "emails",
				MultiValued: true,
				SubAttributes: []schema.CoreAttribute{
					{
						Name: "value",
					},
					{
						Name: "display",
					},
					{
						Name: "type",
						CanonicalValues: []string{
							"work", "home", "other",
						},
					},
					{
						Name: "primary",
						Type: schema.DataTypeBoolean,
					},
				},
			}),
		},
	}

	userSchemaExtension := schema.Schema{
		ID:          "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		Name:        "EnterpriseUser",
		Description: "Enterprise User",
		Attributes: []schema.CoreAttribute{
			{
				Name: "employeeNumber",
			},
			{
				Name: "organization",
			},
		},
	}

	return Server{
		Config: ServiceProviderConfig{},
		ResourceTypes: []ResourceType{
			{
				ID:          "User",
				Name:        "User",
				Endpoint:    "/Users",
				Description: "User Account",
				Schema:      userSchema,
				Handler:     newTestResourceHandler(),
			},
			{
				ID:          "EnterpriseUser",
				Name:        "EnterpriseUser",
				Endpoint:    "/EnterpriseUser",
				Description: "Enterprise User Account",
				Schema:      userSchema,
				SchemaExtensions: []SchemaExtension{
					{Schema: userSchemaExtension},
				},
				Handler: newTestResourceHandler(),
			},
		},
	}
}

func newTestResourceHandler() ResourceHandler {
	data := make(map[string]ResourceAttributes)

	// Generate enough test data to test pagination
	for i := 1; i < 21; i++ {
		data[fmt.Sprintf("000%d", i)] = ResourceAttributes{
			"userName": fmt.Sprintf("test%d", i),
		}
	}

	return testResourceHandler{
		data: data,
	}
}

func TestInvalidEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v2/Invalid", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestServerSchemasEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Schemas", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.TotalResults != 2 {
		t.Errorf("handler returned unexpected body: got %v want 2 total result", rr.Body.String())
	}

	if len(response.Resources.([]interface{})) != 2 {
		t.Fatal("resources contains more than one schema")
	}

	s, ok := response.Resources.([]interface{})[0].(map[string]interface{})
	if !ok {
		t.Fatal("schema is not an object")
	}

	id, ok := s["id"].(string)

	if !ok && id != "urn:ietf:params:scim:schemas:core:2.0:User" &&
		id != "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User" {
		t.Errorf("schema does not contain the correct id: %v", id)
		log.Printf("%#v\n", response.Resources.([]interface{})[0])
	}
}

func TestServerSchemaEndpointInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Schemas/urn:ietf:params:scim:schemas:core:2.0:Group", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

}

func TestServerSchemaEndpointValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Schemas/urn:ietf:params:scim:schemas:core:2.0:User", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var s map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &s); err != nil {
		t.Fatal(err)
	}

	id, ok := s["id"].(string)
	if !ok && id != "urn:ietf:params:scim:schemas:core:2.0:User" {
		t.Errorf("schema does not contain the correct id: %s", id)
	}
}

func TestServerResourceTypesHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ResourceTypes", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.TotalResults != 2 {
		t.Errorf("handler returned unexpected body: got %v want 1 total result", rr.Body.String())
	}

	if len(response.Resources.([]interface{})) != 2 {
		t.Fatal("resources contains more than one schema")
	}

	resourceType, ok := response.Resources.([]interface{})[0].(map[string]interface{})
	if !ok {
		t.Errorf("resource type is not an object")
	}

	name, ok := resourceType["name"].(string)
	if !ok && name != "User" &&
		name != "EnterpriseUser" {
		t.Errorf("schema does not contain the correct id: %v", resourceType["name"])
	}
}

func TestServerResourceTypeHandlerInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ResourceTypes/Group", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestServerResourceTypeHandlerValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ResourceTypes/User", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resourceType map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resourceType); err != nil {
		t.Fatal(err)
	}
	if resourceType["id"] != "User" {
		t.Errorf("schema does not contain the correct name: %s", resourceType["name"])
	}
}

func TestServerServiceProviderConfigHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ServiceProviderConfig", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestServerResourcePostHandlerInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/Users", strings.NewReader(`{"id": "other"}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestServerResourcePostHandlerValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/Users", strings.NewReader(`{"id": "other", "userName": "test1"}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}
	if resource["userName"] != "test1" {
		t.Error("handler did not return the resource correctly")
	}
}

func TestServerResourceGetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Users/0001", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}
	if resource["userName"] != "test1" {
		t.Error("handler did not return the resource correctly")
	}
}

func TestServerResourceGetHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Users/9999", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var scimErr scimError
	if err := json.Unmarshal(rr.Body.Bytes(), &scimErr); err != nil {
		t.Error(err)
	}
	if scimErr != scimErrorResourceNotFound("9999") {
		t.Errorf("wrong scim error: %v", scimErr)
	}
}

func TestServerResourcesGetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Users", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.TotalResults != 20 {
		t.Errorf("handler returned unexpected body: got %v want 20 total result", response.TotalResults)
	}
}

func TestServerResourcesGetHandlerPagination(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Users?count=2&startIndex=2", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.TotalResults != 20 {
		t.Errorf("handler returned unexpected body: got %v want 20 total result", response.TotalResults)
	}
}

func TestServerResourcesGetHandlerMaxCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Users?count=20000", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ListResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}

	if response.TotalResults != 20 {
		t.Errorf("handler returned unexpected body: got %v want 20 total result", response.TotalResults)
	}
}

// Tests valid add, replace, and remove operations
func TestServerResourcePatchHandlerValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/Users/0001", strings.NewReader(`{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
		"Operations":[
		  {
		    "op":"add",
		    "value":{
		      "emails":[
		        {
			  "value":"babs@jensen.org",
			  "type":"home"
		        }
		      ]
		    }
		  },
		  {
		    "op":"replace",
		    "path":"active",
		    "value":false
		  },
		  {
		    "op":"remove",
		    "path":"displayName"
		  }
		]
	}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Error response: %v\n", resource)
	}

	if resource["displayName"] != nil {
		t.Errorf("handler did not remove the displayName attribute")
	}

	if resource["active"] != false {
		t.Errorf("handler did not deactivate user")
	}

	if resource["emails"] == nil || len(resource["emails"].([]interface{})) < 1 {
		t.Errorf("handler did not add user's email address")
	}
}

func TestServerResourcePatchHandlerFailOnBadType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/Users/0001", strings.NewReader(`{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
		"Operations":[
		  {
		    "op":"replace",
		    "path":"active",
		    "value":"test"
		  }
		]
	}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Error response: %v\n", resource)
	}
}

func TestServerResourcePatchHandlerFailOnUndefinedAttribute(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/Users/0001", strings.NewReader(`{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
		"Operations":[
		  {
		    "op":"add",
		    "value":{
		      "notActuallyAThing": "adfad"
		    }
		  }
		]
	}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Error response: %v\n", resource)
	}
}

func runPatchImmutableTest(t *testing.T, op, path string, expectedStatus int) {
	req := httptest.NewRequest(http.MethodPatch, "/Users/0001", strings.NewReader(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
		"Operations":[
		  {
		    "op":"%s",
		    "path":"%s",
		    "value":"test"
		  }
		]
	}`, op, path)))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Error response: %v\n", resource)
	}
}

// Ensure we error when changing an immutable or readonly property while allowing adding of immutable properties.
func TestServerResourcePatchHandlerFailOnImmutable(t *testing.T) {
	runPatchImmutableTest(t, PatchOperationAdd, "immutableThing", http.StatusOK)
	runPatchImmutableTest(t, PatchOperationRemove, "immutableThing", http.StatusBadRequest)
	runPatchImmutableTest(t, PatchOperationReplace, "immutableThing", http.StatusBadRequest)
	runPatchImmutableTest(t, PatchOperationReplace, "readonlyThing", http.StatusBadRequest)
	runPatchImmutableTest(t, PatchOperationRemove, "readonlyThing", http.StatusBadRequest)
	runPatchImmutableTest(t, PatchOperationReplace, "readonlyThing", http.StatusBadRequest)
}

func TestServerResourcePutHandlerInvalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/Users/0001", strings.NewReader(`{"more": "test"}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestServerResourcePutHandlerValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/Users/0001", strings.NewReader(`{"id": "test", "userName": "other"}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resource); err != nil {
		t.Fatal(err)
	}
	if resource["userName"] != "other" {
		t.Errorf("handler did not replace previous resource")
	}
}

func TestServerResourcePutHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/Users/9999", strings.NewReader(`{"userName": "other"}`))
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var scimErr scimError
	if err := json.Unmarshal(rr.Body.Bytes(), &scimErr); err != nil {
		t.Error(err)
	}

	if scimErr != scimErrorResourceNotFound("9999") {
		t.Errorf("wrong scim error: %v", scimErr)
	}
}

func TestServerResourceDeleteHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/Users/0001", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestServerResourceDeleteHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/Users/9999", nil)
	rr := httptest.NewRecorder()
	newTestServer().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var scimErr scimError
	if err := json.Unmarshal(rr.Body.Bytes(), &scimErr); err != nil {
		t.Error(err)
	}

	if scimErr != scimErrorResourceNotFound("9999") {
		t.Errorf("wrong scim error: %v", scimErr)
	}
}
