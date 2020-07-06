package coordinator

import (
	"context"
	"encoding/json"
	goerrors "errors"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/manifoldco/go-base64"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

type gqlResponse struct {
	Data interface{} `json:"data"`
	Err  []string    `json:"errors"`
}

func setupGQLTestServer(resp *gqlResponse) (*httptest.Server, *primitives.URL, error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(resp)
	}))

	clientURL, err := primitives.NewURL(ts.URL)
	if err != nil {
		return nil, nil, err
	}

	return ts, clientURL, nil
}

func createHTTPTransport(ts *httptest.Server, clientURL *primitives.URL, client *http.Client) *HTTPTransport {
	if client == nil {
		client = ts.Client()
	}

	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(client))
	ct := &HTTPTransport{
		client:     gql,
		authToken:  base64.New([]byte("faketoken")),
		url:        clientURL,
		httpClient: client,
	}

	return ct
}

func TestHTTPTransportRaw(t *testing.T) {
	ts, clientURL, err := setupGQLTestServer(&gqlResponse{
		Data: "blahblah",
	})
	if err != nil {
		t.Errorf("error setting up test")
	}
	defer ts.Close()

	t.Run("no variables, no errors", func(t *testing.T) {
		ct := createHTTPTransport(ts, clientURL, nil)
		err := ct.Raw(context.TODO(), "fakequery", nil, nil)
		if err != nil {
			t.Errorf("failed call to Raw with error %v", err)
		}
	})

	t.Run("with variables, no errors", func(t *testing.T) {
		ct := createHTTPTransport(ts, clientURL, nil)
		variables := map[string]interface{}{
			"one": 1,
			"two": 2,
		}
		err := ct.Raw(context.TODO(), "fakequery", variables, nil)
		if err != nil {
			t.Errorf("failed call to Raw with error %v", err)
		}
	})

	t.Run("fails on net error", func(t *testing.T) {
		hc := &http.Client{
			Transport: ErrorRoundTripper{
				err: &net.DNSConfigError{
					Err: goerrors.New("generic error"),
				},
			},
		}

		ct := createHTTPTransport(ts, clientURL, hc)
		variables := map[string]interface{}{
			"one": 1,
			"two": 2,
		}
		err = ct.Raw(context.TODO(), "fakequery", variables, nil)
		if !errors.FromCause(err, NetworkCause) {
			t.Errorf("incorrect error returned: %v", err)
		}
	})
}

func TestHTTPTransportURL(t *testing.T) {
	ts, clientURL, err := setupGQLTestServer(&gqlResponse{
		Data: "blahblah",
	})
	if err != nil {
		t.Errorf("error setting up test")
	}
	defer ts.Close()

	ct := createHTTPTransport(ts, clientURL, nil)
	if ct.URL() != clientURL {
		t.Errorf("did not return correct url: got %q, want %q", ct.URL(), clientURL)
	}
}

func TestHTTPTransportTokenLogin(t *testing.T) {
	wantSess := primitives.Session{}
	ts, clientURL, err := setupGQLTestServer(&gqlResponse{
		Data: &SessionResponse{
			Session: wantSess,
		},
	})
	if err != nil {
		t.Errorf("error setting up test")
	}
	defer ts.Close()

	ct := createHTTPTransport(ts, clientURL, nil)
	gotSess, err := ct.TokenLogin(context.TODO(), &auth.APIToken{})
	if err != nil {
		t.Errorf("token login returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(gotSess, &wantSess) {
		t.Errorf("received bad session: got %v want %v", gotSess, wantSess)
	}

	// TODO(thor): also test error cases
}

func TestHTTPTransportEmailLogin(t *testing.T) {
	wantSess := primitives.Session{}
	ts, clientURL, err := setupGQLTestServer(&gqlResponse{
		Data: &SessionResponse{
			Session: wantSess,
		},
	})
	if err != nil {
		t.Errorf("error setting up test")
	}
	defer ts.Close()

	ct := createHTTPTransport(ts, clientURL, nil)
	gotSess, err := ct.EmailLogin(context.TODO(), primitives.Email{}, primitives.EmptyPassword)
	if err != nil {
		t.Errorf("email login returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(gotSess, &wantSess) {
		t.Errorf("received bad session: got %v want %v", gotSess, wantSess)
	}

	// TODO(thor): also test error cases
}

func TestHTTPTransportLogout(t *testing.T) {
	ts, clientURL, err := setupGQLTestServer(&gqlResponse{})
	if err != nil {
		t.Errorf("error setting up test")
	}
	defer ts.Close()

	ct := createHTTPTransport(ts, clientURL, nil)
	err = ct.Logout(context.TODO(), nil)
	if err != nil {
		t.Errorf("logout returned unexpected error: %v", err)
	}

	// TODO(thor): also test error cases
}

type ErrorRoundTripper struct {
	err error
}

func (e ErrorRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, e.err
}
