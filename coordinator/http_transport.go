package coordinator

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/machinebox/graphql"
	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

// HTTPTransport is a ClientTransport that interacts with the Coordinator via
// GraphQL over HTTP.
type HTTPTransport struct {
	client     *graphql.Client
	authToken  *base64.Value
	url        *primitives.URL
	httpClient *http.Client
}

// NewHTTPTransport returns a ClientTransport configured to make requests via
// GraphQL over HTTP
func NewHTTPTransport(coordinatorURL *primitives.URL, authToken *base64.Value, certFile string) ClientTransport {
	// this function never actually returns an error...
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}
	if certFile != "" {
		httpClient = httpsClient(certFile)
	}

	if authToken != nil && len(authToken.String()) == 0 {
		authToken = nil
	}

	client := graphql.NewClient(coordinatorURL.String()+"/v1/query", graphql.WithHTTPClient(httpClient))
	transport := &HTTPTransport{
		client:     client,
		httpClient: httpClient,
		url:        coordinatorURL,
	}

	transport.SetToken(authToken)

	return transport
}

// URL returns the underlying URL used by this Transport
func (c *HTTPTransport) URL() *primitives.URL {
	return c.url
}

// Raw wraps the NewRequest and does common req changes like adding authorization
// headers. It calls Run passing the object to be filled with the request data.
func (c *HTTPTransport) Raw(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
	req := graphql.NewRequest(query)

	for key, val := range variables {
		req.Var(key, val)
	}

	err := c.client.Run(ctx, req, resp)
	if err != nil {
		if nerr, ok := err.(net.Error); ok {
			return errors.New(NetworkCause, "Could not contact coordinator: %s", nerr.Error())
		}
		return convertError(err)
	}

	return nil
}

type SessionResponse struct {
	Session primitives.Session `json:"createSession"`
}

// TokenLogin enables a user or service to login using an APIToken
func (c *HTTPTransport) TokenLogin(ctx context.Context, apiToken *auth.APIToken) (*primitives.Session, error) {
	return tokenLogin(ctx, c, apiToken)
}

// EmailLogin starts step 1 of the login flow using an email & password
func (c *HTTPTransport) EmailLogin(ctx context.Context, email primitives.Email, password primitives.Password) (*primitives.Session, error) {
	return emailLogin(ctx, c, email, password)
}

// Logout of the active session
func (c *HTTPTransport) Logout(ctx context.Context, authToken *base64.Value) error {
	return logout(ctx, c, authToken)
}

// Authenticated returns whether the client is authenticated or not.
// If the authToken is not nil then its authenticated!
func (c *HTTPTransport) Authenticated() bool {
	return c.authToken != nil
}

// SetToken enables a caller to set the auth token used by the transport
func (c *HTTPTransport) SetToken(value *base64.Value) {
	c.authToken = value

	if c.authToken != nil {
		cookie := &http.Cookie{
			Name:  "token",
			Value: c.authToken.String(),
		}
		c.httpClient.Jar.SetCookies(c.url.URL, []*http.Cookie{cookie})
	}
}

// Token enables a caller to retrieve the current auth token used by the
// transport
func (c *HTTPTransport) Token() *base64.Value {
	return c.authToken
}

func convertError(err error) error {
	return errors.New(errors.UnknownCause, strings.TrimPrefix(err.Error(), "graphql: "))
}

func httpsClient(certFile string) *http.Client {
	caCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{RootCAs: caCertPool}

	// The below configuration is borrowed from the DefaultTransport
	//
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}

	return &http.Client{Transport: tr}
}

// Post does a raw http POST to the specified url
func (c *HTTPTransport) Post(url string, req interface{}) ([]byte, error) {
	contentType := "application/json; charset=utf-8"

	by, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Post(url, contentType, bytes.NewBuffer(by))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		e := &errors.Error{}
		err := json.Unmarshal(body, e)
		if err != nil {
			return nil, err
		}

		return nil, e
	}

	return body, nil
}
