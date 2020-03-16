package auth

import (
	"crypto/ed25519"
	"time"

	"github.com/manifoldco/go-base64"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

// TokenType is an enum holding the category of sessions
type TokenType string

var (
	// Login is the session type used during the login flow
	Login TokenType = "login"
	// Authenticated is the session type used on normal API calls
	Authenticated TokenType = "auth"
)

//
var (
	// LoginTokenDuration how long login tokens last
	LoginTokenDuration = time.Minute * 5

	// AuthTokenDuration how long auth tokens last
	AuthTokenDuration = time.Hour * 24
)

// TokenAuthority is the authority over token, it generates
// and verifies tokens based on the private/public keys it owns
type TokenAuthority struct {
	privateKey   ed25519.PrivateKey
	PublicKey    ed25519.PublicKey
	ServiceEmail string
}

// NewTokenAuthority returns a new token authority
func NewTokenAuthority(serviceEmail string) (*TokenAuthority, error) {
	pub, priv, err := newKey()
	if err != nil {
		return nil, err
	}

	return &TokenAuthority{
		privateKey:   priv,
		PublicKey:    pub,
		ServiceEmail: serviceEmail,
	}, nil
}

// Verify verifies that a JWT token was signed by the correct private key
func (t *TokenAuthority) Verify(signedToken *base64.Value) error {
	tok, err := jwt.ParseSigned(string(*signedToken))
	if err != nil {
		return err
	}

	claims := jwt.Claims{}
	err = tok.Claims(t.PublicKey, &claims)
	if err != nil {
		return err
	}

	return claims.Validate(jwt.Expected{
		Issuer: t.ServiceEmail,
		Time:   time.Now().UTC(), // time used to compare expiry and not before
	})
}

// Generate generates a JWT with 4 claims:
// - Expiry: time the JWT expires recommend 5 minutes for login sessions
//              and 24 hours for general authenticated sessions
// - IssuedAt: time the JWT was issued
// - NotBefore: the JWT will not be accepted before this time has passed
// - Issuer: the service email of the issuing controller
func (t *TokenAuthority) Generate(tokenType TokenType) (*base64.Value, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: t.privateKey},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	var expiresIn time.Time
	switch tokenType {
	case Login:
		expiresIn = now.Add(LoginTokenDuration)
	case Authenticated:
		expiresIn = now.Add(AuthTokenDuration)
	default:
		return nil, errors.New(InvalidTokenType, "Invalid token type must be login or authenticated")
	}

	cl := jwt.Claims{
		Issuer:    t.ServiceEmail,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Expiry:    jwt.NewNumericDate(expiresIn),
	}

	signedToken, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return nil, err
	}

	return base64.New([]byte(signedToken)), nil
}
