package coordinator

import (
	"context"
	"encoding/json"
	goerrors "errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/machinebox/graphql"
	"github.com/manifoldco/go-base64"
)

func TestClientTransportRaw(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{ \"data\": \"blahblah\", \"errors\": [] }")
	}))
	defer ts.Close()

	clientURL, err := primitives.NewURL(ts.URL)
	if err != nil {
		t.Errorf("error creating test client")
	}

	t.Run("no variables, no errors", func(t *testing.T) {
		gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(ts.Client()))
		ct := &ClientTransport{
			client:    gql,
			authToken: base64.New([]byte("faketoken")),
			url:       clientURL,
		}

		err = ct.Raw(context.TODO(), "fakequery", nil, nil)
		if err != nil {
			t.Errorf("failed call to Raw with error %v", err)
		}
	})

	t.Run("with variables, no errors", func(t *testing.T) {
		gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(ts.Client()))
		ct := &ClientTransport{
			client:    gql,
			authToken: base64.New([]byte("faketoken")),
			url:       clientURL,
		}

		variables := map[string]interface{}{
			"one": 1,
			"two": 2,
		}
		err = ct.Raw(context.TODO(), "fakequery", variables, nil)
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
		gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(hc))
		ct := &ClientTransport{
			client:    gql,
			authToken: base64.New([]byte("faketoken")),
			url:       clientURL,
		}

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

func TestClientTransportURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{ \"data\": \"blahblah\", \"errors\": [] }")
	}))
	defer ts.Close()

	clientURL, err := primitives.NewURL(ts.URL)
	if err != nil {
		t.Errorf("error creating test client")
	}

	gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(ts.Client()))
	ct := &ClientTransport{
		client:    gql,
		authToken: base64.New([]byte("faketoken")),
		url:       clientURL,
	}

	if ct.URL() != clientURL {
		t.Errorf("did not return correct url: got %q, want %q", ct.URL(), clientURL)
	}
}

func TestClientTransportTokenLogin(t *testing.T) {
	var resp struct {
		Data SessionResponse `json:"data"`
		Err  []string        `json:"errors"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	clientURL, err := primitives.NewURL(ts.URL)
	if err != nil {
		t.Errorf("error creating test client")
	}

	gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(ts.Client()))
	ct := &ClientTransport{
		client:    gql,
		authToken: base64.New([]byte("faketoken")),
		url:       clientURL,
	}

	wantSess := primitives.Session{}
	resp.Data = SessionResponse{
		Session: wantSess,
	}
	gotSess, err := ct.TokenLogin(context.TODO(), &auth.APIToken{})
	if err != nil {
		t.Errorf("token login returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(gotSess, &wantSess) {
		t.Errorf("received bad session: got %v want %v", gotSess, wantSess)
	}

	// TODO(thor): also test error cases
}

func TestClientTransportEmailLogin(t *testing.T) {
	var resp struct {
		Data SessionResponse `json:"data"`
		Err  []string        `json:"errors"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	clientURL, err := primitives.NewURL(ts.URL)
	if err != nil {
		t.Errorf("error creating test client")
	}

	gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(ts.Client()))
	ct := &ClientTransport{
		client:    gql,
		authToken: base64.New([]byte("faketoken")),
		url:       clientURL,
	}

	wantSess := primitives.Session{}
	resp.Data = SessionResponse{
		Session: wantSess,
	}
	gotSess, err := ct.EmailLogin(context.TODO(), primitives.Email{}, primitives.EmptyPassword)
	if err != nil {
		t.Errorf("email login returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(gotSess, &wantSess) {
		t.Errorf("received bad session: got %v want %v", gotSess, wantSess)
	}

	// TODO(thor): also test error cases
}

func TestClientTransportLogout(t *testing.T) {
	var resp struct {
		Data string   `json:"data"`
		Err  []string `json:"errors"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	clientURL, err := primitives.NewURL(ts.URL)
	if err != nil {
		t.Errorf("error creating test client")
	}

	gql := graphql.NewClient(ts.URL+"/v1/query", graphql.WithHTTPClient(ts.Client()))
	ct := &ClientTransport{
		client:    gql,
		authToken: base64.New([]byte("faketoken")),
		url:       clientURL,
	}

	resp.Data = ""
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
