package go3scale

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
)

const (
	providerKey = "providerKey"
	hostKey     = "example.com"
	userKey     = "userKey"
)

var authTests = []struct {
	status   int
	expected bool
}{
	{http.StatusOK, true}, {http.StatusConflict, true}, {http.StatusBadRequest, false}, {http.StatusForbidden, false},
}

var usages = []Usage{
	{"Test", 1}, {"Another", 2},
}

func TestCreateURL(t *testing.T) {
	expected := "https://example.com/transactions/authrep.xml?user_key=userKey&provider_key=providerKey&usage[Test]=1&usage[Another]=2"
	c := NewClient(providerKey, hostKey)

	URL := createURL(c, "userKey", usages)

	if URL != expected {
		t.Errorf("Received '%v' but expected '%v'", URL, expected)
	}
}

// TestAuthrepUserKey loops through our authTests array so that we can test a
// variety of test cases.
func TestAuthrepUserKey(t *testing.T) {
	for _, tt := range authTests {
		server, client := testAPICall(t, tt.status, providerKey, hostKey)
		defer server.Close()

		result, err := client.AuthrepUserKey(userKey, usages)

		if result.IsSuccess() != tt.expected {
			t.Errorf("Expected IsSuccess to be %v. Code recieved: %v", tt.expected, result.Code)
		}

		if err != nil {
			t.Errorf("Received an error: %v", err)
		}
	}
}

// RewriteTransport is an http.RoundTripper that rewrites requests
// using the provided URL's Scheme and Host, and its Path as a prefix.
// The Opaque field is untouched.
// If Transport is nil, http.DefaultTransport is used
type RewriteTransport struct {
	Transport http.RoundTripper
	URL       *url.URL
}

func (t RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// note that url.URL.ResolveReference doesn't work here
	// since t.u is an absolute url
	req.URL.Scheme = t.URL.Scheme
	req.URL.Host = t.URL.Host
	req.URL.Path = path.Join(t.URL.Path, req.URL.Path)
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(req)
}

func testAPICall(t *testing.T, code int, providerKey, host string) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(userAgentHeaderKey) != userAgentHeaderValue {
			t.Errorf("%v did not contain %v", userAgentHeaderKey, userAgentHeaderValue)
		}
		w.WriteHeader(code)
		fmt.Fprintln(w, "")
	}))

	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatalln("failed to parse httptest.Server URL:", err)
	}

	transport := &RewriteTransport{URL: u}

	httpClient := &http.Client{Transport: transport}
	client := &Client{
		ProviderKey: providerKey,
		Host:        host,
		httpClient:  httpClient,
	}

	return server, client
}
