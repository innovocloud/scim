package scim

import (
	"encoding/json"
)

// ServiceProviderConfig enables a service provider to discover SCIM specification features in a standardized form as
// well as provide additional implementation details to clients.
type ServiceProviderConfig struct {
	// DocumentationURI is an HTTP-addressable URL pointing to the service provider's human-consumable help
	// documentation.
	DocumentationURI string `json:"documentationUri,omitempty"`
	// AuthenticationSchemes is a multi-valued complex type that specifies supported authentication scheme properties.
	AuthenticationSchemes []AuthenticationScheme `json:",omitempty"`
	// MaxResults denotes the the integer value specifying the maximum number of resources returned in a response. It defaults to 100.
	MaxResults int
	// SupportFiltering whether you SCIM implementation will support filtering.
	SupportFiltering bool
	// SupportPatch whether your SCIM implementation will support patch requests.
	SupportPatch          bool
	SupportChangePassword bool
	SupportSort           bool
	SupportEtag           bool
	SupportBulk           bool
	BulkMaxOpts           int
	BulkMaxPayload        int
	// OmitSchemasHeader is set to true to not return the list metadata on the /Schemas endpoint for compatibility with OneIdentity products
	OmitSchemasHeader bool
	// OmitResourceTypesHeader is set to true to not return the list metadata on the /ResourceTypes endpoint for compatibility with OneIdentity products
	OmitResourceTypesHeader bool
}

// AuthenticationScheme specifies a supported authentication scheme property.
type AuthenticationScheme struct {
	// Type is the authentication scheme. This specification defines the values "oauth", "oauth2", "oauthbearertoken",
	// "httpbasic", and "httpdigest".
	Type AuthenticationType `json:"type"`
	// Name is the common authentication scheme name, e.g., HTTP Basic.
	Name string `json:"name"`
	// Description of the authentication scheme.
	Description string `json:"description"`
	// SpecURI is an HTTP-addressable URL pointing to the authentication scheme's specification.
	SpecURI string `json:"specUri,omitempty"`
	// DocumentationURI is an HTTP-addressable URL pointing to the authentication scheme's usage documentation.
	DocumentationURI string `json:",omitempty"`
	// Primary is a boolean value indicating the 'primary' or preferred authentication scheme.
	Primary bool `json:"primary"`
}

// AuthenticationType is a single keyword indicating the authentication type of the authentication scheme.
type AuthenticationType string

const (
	// AuthenticationTypeOauth indicates that the authentication type is OAuth.
	AuthenticationTypeOauth AuthenticationType = "oauth"
	// AuthenticationTypeOauth2 indicates that the authentication type is OAuth2.
	AuthenticationTypeOauth2 AuthenticationType = "oauth2"
	// AuthenticationTypeOauthBearerToken indicates that the authentication type is OAuth2 Bearer Token.
	AuthenticationTypeOauthBearerToken AuthenticationType = "oauthbearertoken"
	// AuthenticationTypeHTTPBasic indicated that the authentication type is Basic Access Authentication.
	AuthenticationTypeHTTPBasic AuthenticationType = "httpbasic"
	// AuthenticationTypeHTTPDigest indicated that the authentication type is Digest Access Authentication.
	AuthenticationTypeHTTPDigest AuthenticationType = "httpdigest"
)

// MarshalJSON implements the pkg/encoding/json/Marshaller interface for ServiceProviderConfig
func (config ServiceProviderConfig) MarshalJSON() ([]byte, error) {
	marshalled := map[string]interface{}{
		"schemas":          []string{"urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"},
		"documentationUri": config.DocumentationURI,
		"patch": map[string]bool{
			"supported": config.SupportPatch,
		},
		"bulk": map[string]interface{}{
			"supported":      config.SupportBulk,
			"maxOperations":  config.BulkMaxOpts,
			"maxPayloadSize": config.BulkMaxPayload,
		},
		"filter": map[string]interface{}{
			"supported":  config.SupportFiltering,
			"maxResults": config.MaxResults,
		},
		"changePassword": map[string]bool{
			"supported": config.SupportChangePassword,
		},
		"sort": map[string]bool{
			"supported": config.SupportSort,
		},
		"etag": map[string]bool{
			"supported": config.SupportEtag,
		},
	}

	if config.AuthenticationSchemes != nil {
		marshalled["authenticationSchemes"] = config.AuthenticationSchemes
	}

	return json.Marshal(marshalled)
}

// getItemsPerPage retrieves the configured default count. It falls back to 100 when not configured.
func (config ServiceProviderConfig) getItemsPerPage() int {
	if config.MaxResults < 1 {
		return fallbackCount
	}
	return config.MaxResults
}
