package auth

import (
	"regexp"

	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/manifoldco/go-base64"
)

// ErrorInvalidAuthHeader is an invalid auth header error
var ErrorInvalidAuthHeader = errors.New(InvalidAuthHeader, "Unable to parse auth header")

// GetBearerToken parses a string looking for a bearer token
func GetBearerToken(bearer string) (*base64.Value, error) {
	r := regexp.MustCompile(`Bearer\s(?P<bearer_token>[^\s]+)$`)
	authHeaderParts := r.FindStringSubmatch(bearer)
	if len(authHeaderParts) != 2 {
		return nil, ErrorInvalidAuthHeader
	}

	token, err := base64.NewFromString(authHeaderParts[1])
	if err != nil {
		return nil, ErrorInvalidAuthHeader
	}

	return token, nil
}
