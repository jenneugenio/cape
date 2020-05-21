package auth

import (
	"crypto/ed25519"
	"time"

	"github.com/manifoldco/go-base64"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

var (
	// LoginTokenDuration how long login tokens last
	LoginTokenDuration = time.Minute * 5

	// AuthTokenDuration how long auth tokens last
	AuthTokenDuration = time.Hour * 24
)

// TokenAuthority is the authority over token, it generates
// and verifies tokens based on the private/public keys it owns
type TokenAuthority struct {
	keypair      *Keypair
	serviceEmail string
}

// NewTokenAuthority returns a new token authority
func NewTokenAuthority(keypair *Keypair, serviceEmail string) (*TokenAuthority, error) {
	return &TokenAuthority{
		keypair:      keypair,
		serviceEmail: serviceEmail,
	}, nil
}

// Verify verifies that a JWT token was signed by the correct private key. Returns
// the session ID contained inside of the token.
func (t *TokenAuthority) Verify(signedToken *base64.Value) (database.ID, error) {
	if t.keypair == nil {
		return database.EmptyID, errors.New(MissingKeyPair, "Missing key pair cannot verify token")
	}

	tok, err := jwt.ParseSigned(string(*signedToken))
	if err != nil {
		return database.EmptyID, err
	}

	claims := jwt.Claims{}
	err = tok.Claims(t.PublicKey(), &claims)
	if err != nil {
		return database.EmptyID, err
	}

	err = claims.Validate(jwt.Expected{
		Issuer: t.serviceEmail,
		Time:   time.Now().UTC(), // time used to compare expiry and not before
	})
	if err != nil {
		return database.EmptyID, err
	}

	return database.DecodeFromString(claims.ID)
}

// PublicKey returns a copy of the ed25519 PublicKey
func (t *TokenAuthority) PublicKey() ed25519.PublicKey {
	return t.keypair.PublicKey
}

// Generate generates a JWT with 4 claims:
// - Expiry: time the JWT expires recommend 5 minutes for login sessions
//           and 24 hours for general authenticated sessions
// - IssuedAt: time the JWT was issued
// - NotBefore: the JWT will not be accepted before this time has passed
// - Issuer: the service email of the issuing coordinator
func (t *TokenAuthority) Generate(tokenType primitives.TokenType, sessionID database.ID) (*base64.Value, time.Time, error) {
	if t.keypair == nil {
		return nil, time.Time{}, errors.New(MissingKeyPair, "Missing key pair cannot generate token")
	}

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: t.keypair.PrivateKey},
		(&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return nil, time.Time{}, err
	}

	now := time.Now().UTC()

	var expiresIn time.Time
	switch tokenType {
	case primitives.Login:
		expiresIn = now.Add(LoginTokenDuration)
	case primitives.Authenticated:
		expiresIn = now.Add(AuthTokenDuration)
	default:
		return nil, time.Time{}, errors.New(primitives.InvalidTokenType,
			"Invalid token type must be login or authenticated")
	}

	cl := jwt.Claims{
		ID:        sessionID.String(),
		Issuer:    t.serviceEmail,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Expiry:    jwt.NewNumericDate(expiresIn),
	}

	signedToken, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return nil, time.Time{}, err
	}

	return base64.New([]byte(signedToken)), expiresIn, nil
}

func (t *TokenAuthority) SetKeyPair(keypair *Keypair) {
	t.keypair = keypair
}
