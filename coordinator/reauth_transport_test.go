package coordinator

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func createReAuthTransport(token *auth.APIToken, mock *MockClientTransport) (*ReAuthTransport, error) {
	t, err := NewReAuthTransport(mock.URL(), token)
	if err != nil {
		return nil, err
	}

	reauth := t.(*ReAuthTransport)
	reauth.transport = mock
	return reauth, nil
}

func TestReAuthTransport(t *testing.T) {
	gm.RegisterTestingT(t)

	validCoordURL, err := primitives.NewURL("http://localhost:8080")
	gm.Expect(err).To(gm.BeNil())

	validAPIToken, err := auth.GenerateToken()
	gm.Expect(err).To(gm.BeNil())

	badCoordURL, err := primitives.NewURL("http://localhost:8080")
	gm.Expect(err).To(gm.BeNil())

	badCoordURL.Scheme = "ftp" // ya right?!

	badAPIToken, err := auth.GenerateToken()
	gm.Expect(err).To(gm.BeNil())
	badAPIToken.Version = 0x00

	token, err := auth.GenerateToken()
	gm.Expect(err).To(gm.BeNil())

	_, user, err := primitives.GenerateUser("hello", "hello@everyone.com")
	gm.Expect(err).To(gm.BeNil())

	session, err := primitives.NewSession(user)
	gm.Expect(err).To(gm.BeNil())

	url, err := primitives.NewURL("http://localhost:8080")
	gm.Expect(err).To(gm.BeNil())

	t.Run("cannot pass invalid data", func(t *testing.T) {
		tests := []struct {
			name           string
			coordinatorURL *primitives.URL
			apiToken       *auth.APIToken
			cause          *errors.Cause
		}{
			{
				name:           "valid input",
				coordinatorURL: validCoordURL,
				apiToken:       validAPIToken,
			},
			{
				name:     "missing coord url",
				apiToken: validAPIToken,
				cause:    &InvalidArgumentCause,
			},
			{
				name:           "missing api token",
				coordinatorURL: validCoordURL,
				cause:          &InvalidArgumentCause,
			},
			{
				name:           "bad coord url",
				coordinatorURL: badCoordURL,
				apiToken:       validAPIToken,
				cause:          &primitives.InvalidURLCause,
			},
			{
				name:           "bad api token",
				coordinatorURL: validCoordURL,
				apiToken:       badAPIToken,
				cause:          &auth.BadAPITokenVersion,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewReAuthTransport(tc.coordinatorURL, tc.apiToken)
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})

	t.Run("authenticates first if not authenticated", func(t *testing.T) {
		mock, err := NewMockClientTransport(url, []*MockResponse{
			{
				Value: SessionResponse{
					Session: *session,
				},
			},
			{
				Value: MeResponse{
					Identity: user.IdentityImpl,
				},
			},
		})
		gm.Expect(err).To(gm.BeNil())

		reauth, err := createReAuthTransport(token, mock)
		gm.Expect(err).To(gm.BeNil())

		client := &Client{transport: reauth}
		identity, err := client.Me(context.TODO())
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(mock.Requests)).To(gm.Equal(2))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("CreateSession"))
		gm.Expect(mock.Requests[0].Variables).To(gm.HaveKey("token_id"))
		gm.Expect(mock.Requests[0].Variables).To(gm.HaveKey("secret"))
		gm.Expect(mock.Requests[1].Query).To(gm.ContainSubstring("Me"))
		gm.Expect(identity.GetEmail()).To(gm.Equal(user.GetEmail()))
	})

	t.Run("authenticates again if ErrAuthentication received after first auth", func(t *testing.T) {
		mock, err := NewMockClientTransport(url, []*MockResponse{
			{
				Value: SessionResponse{
					Session: *session,
				},
			},
			{
				Error: auth.ErrAuthentication,
			},
			{
				Value: SessionResponse{
					Session: *session,
				},
			},
			{
				Value: MeResponse{
					Identity: user.IdentityImpl,
				},
			},
		})
		gm.Expect(err).To(gm.BeNil())

		reauth, err := createReAuthTransport(token, mock)
		gm.Expect(err).To(gm.BeNil())

		client := &Client{transport: reauth}
		identity, err := client.Me(context.TODO())
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(mock.Requests)).To(gm.Equal(4))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("CreateSession"))
		gm.Expect(mock.Requests[0].Variables).To(gm.HaveKey("token_id"))
		gm.Expect(mock.Requests[0].Variables).To(gm.HaveKey("secret"))
		gm.Expect(mock.Requests[1].Query).To(gm.ContainSubstring("Me"))
		gm.Expect(mock.Requests[2].Query).To(gm.ContainSubstring("CreateSession"))
		gm.Expect(mock.Requests[2].Variables).To(gm.HaveKey("token_id"))
		gm.Expect(mock.Requests[2].Variables).To(gm.HaveKey("secret"))
		gm.Expect(mock.Requests[3].Query).To(gm.ContainSubstring("Me"))
		gm.Expect(identity.GetEmail()).To(gm.Equal(user.GetEmail()))
	})

	t.Run("returns error if authentication fails", func(t *testing.T) {
		mock, err := NewMockClientTransport(url, []*MockResponse{
			{
				Error: auth.ErrAuthentication,
			},
		})
		gm.Expect(err).To(gm.BeNil())

		reauth, err := createReAuthTransport(token, mock)
		gm.Expect(err).To(gm.BeNil())

		client := &Client{transport: reauth}
		_, err = client.Me(context.TODO())
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("returns error if second auth fails", func(t *testing.T) {
		mock, err := NewMockClientTransport(url, []*MockResponse{
			{
				Value: SessionResponse{
					Session: *session,
				},
			},
			{
				Error: auth.ErrAuthentication,
			},
			{
				Error: auth.ErrAuthentication,
			},
		})
		gm.Expect(err).To(gm.BeNil())

		reauth, err := createReAuthTransport(token, mock)
		gm.Expect(err).To(gm.BeNil())

		client := &Client{transport: reauth}
		_, err = client.Me(context.TODO())
		gm.Expect(err).ToNot(gm.BeNil())

		gm.Expect(len(mock.Requests)).To(gm.Equal(3))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("CreateSession"))
		gm.Expect(mock.Requests[0].Variables).To(gm.HaveKey("token_id"))
		gm.Expect(mock.Requests[0].Variables).To(gm.HaveKey("secret"))
		gm.Expect(mock.Requests[1].Query).To(gm.ContainSubstring("Me"))
		gm.Expect(mock.Requests[2].Query).To(gm.ContainSubstring("CreateSession"))
		gm.Expect(mock.Requests[2].Variables).To(gm.HaveKey("token_id"))
		gm.Expect(mock.Requests[2].Variables).To(gm.HaveKey("secret"))
	})
}
