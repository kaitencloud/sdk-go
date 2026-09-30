// Package sdk provides typed Go clients for the Kaiten API.
//
// Kaiten publishes two contracts on two listeners, and this package mirrors that split
// rather than hiding it:
//
//   - Client covers the Core API -- customers, instances, entitlements, licenses,
//     feature flags -- authenticated by an organization-scoped token (`ksh_...`).
//   - PlatformClient covers the Platform API -- organizations, users, and the
//     credentials themselves -- authenticated by a platform credential (`ksm_...`).
//
// The two are not interchangeable in either direction: the Core listener refuses a
// platform credential outright, and the Platform listener serves none of the Core
// paths. A process that needs both builds both, against different base URLs.
package sdk

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kaitencloud/sdk-go/internal/gen"
)

// DefaultTimeout bounds every request made by the HTTP client NewClient builds when the
// caller supplies none. WithHTTPClient replaces that client, and with it this timeout.
const DefaultTimeout = 30 * time.Second

// Client is a typed client for the Kaiten Core API -- the listener a tenant's own
// organization-scoped token talks to. Construct one with NewClient.
//
// Organizations, users and platform credentials are not here: they belong to the
// Platform API, which is a separate listener with a separate credential class, and
// PlatformClient is what reaches it.
//
// The zero value is not usable: the resource namespaces below are nil until NewClient
// populates them, so calling a method on one panics. There is no other way to build a Client.
type Client struct {
	raw *gen.ClientWithResponses

	// list carries the plumbing the paginated list methods reuse. They build
	// their own requests because internal/gen predates Core's pagination and
	// cannot send a cursor, so they need the base URL, HTTP client and
	// request editors the generated client was handed.
	list listTransport

	Components        *Components
	Customers         *Customers
	DeploymentZones   *DeploymentZones
	EntitlementGroups *EntitlementGroups
	Entitlements      *Entitlements
	FeatureFlags      *FeatureFlags
	Instances         *Instances
	LicenseFamilies   *LicenseFamilies
	Licenses          *Licenses
	MetadataFields    *MetadataFields
	Releases          *Releases
	ServiceAccounts   *ServiceAccounts
}

// NewClient creates a Client for the Kaiten Core API at baseURL, applying the given
// Options. baseURL includes the API's path prefix, e.g. https://kaiten.example.com/api.
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	resolved := resolveOptions(baseURL, opts)
	if strings.TrimSpace(resolved.baseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	raw, err := gen.NewClientWithResponses(resolved.baseURL, resolved.coreOptions()...)
	if err != nil {
		return nil, fmt.Errorf("create SDK client: %w", err)
	}

	return newClient(raw, resolved.listTransport()), nil
}

func newClient(raw *gen.ClientWithResponses, list listTransport) *Client {
	client := &Client{
		raw:  raw,
		list: list,
	}

	client.Components = &Components{client: client}
	client.Customers = &Customers{client: client}
	client.DeploymentZones = &DeploymentZones{client: client}
	client.EntitlementGroups = &EntitlementGroups{client: client}
	client.Entitlements = &Entitlements{client: client}
	client.FeatureFlags = &FeatureFlags{client: client}
	client.Instances = &Instances{client: client}
	client.LicenseFamilies = &LicenseFamilies{client: client}
	client.Licenses = &Licenses{client: client}
	client.MetadataFields = &MetadataFields{client: client}
	client.Releases = &Releases{client: client}
	client.ServiceAccounts = &ServiceAccounts{client: client}

	return client
}

// listTransport is what a hand-built request needs to be indistinguishable
// from a generated one: the same base URL, the same HTTP client, and the same
// request editors -- which is where WithBearerToken lives, so a request that
// skipped them would be unauthenticated.
type listTransport struct {
	baseURL string
	doer    *http.Client
	editors []RequestEditorFn
}

// resolve joins an operation path and query onto the base URL exactly as the
// generated requests do: a trailing slash on the server, the path made
// relative, and the two resolved against each other. Anything else would
// drop the "/api" prefix every base URL carries.
func (t listTransport) resolve(path string, query url.Values) (string, error) {
	server := t.baseURL
	if !strings.HasSuffix(server, "/") {
		server += "/"
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, "/") {
		path = "." + path
	}

	resolved, err := serverURL.Parse(path)
	if err != nil {
		return "", err
	}
	resolved.RawQuery = query.Encode()

	return resolved.String(), nil
}

// Raw returns the underlying generated Core API client for low-level access.
func (c *Client) Raw() *gen.ClientWithResponses {
	return c.raw
}
