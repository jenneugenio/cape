package capepg

import (
	"context"
	"github.com/capeprivacy/cape/coordinator/db"
	"github.com/capeprivacy/cape/models"
	"time"
)

type pgToken struct {
	pool    Pool
	timeout time.Duration
}

var _ db.TokensDB = &pgToken{}

// models.Token annotates the credentials field with `json:"-"` to ensure
// that credentials are not sent over the wire.
//
// However, we do want to store creds in the database, so we wrap the model
// token with one that is more friendly for db writes
type dbToken struct {
	*models.Token
	Credentials *models.Credentials `json:"credentials"`
}

func (p pgToken) Get(ctx context.Context, ID string) (*models.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var token dbToken
	s := "select data from tokens where id = $1;"
	row := p.pool.QueryRow(ctx, s, ID)
	err := row.Scan(&token)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		ID:          token.ID,
		UserID:      token.UserID,
		Credentials: token.Credentials,
	}, nil
}

func (p pgToken) Create(ctx context.Context, token models.Token) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	t := dbToken{
		Token:       &token,
		Credentials: token.Credentials,
	}

	s := "insert into tokens (data) values ($1);"
	_, err := p.pool.Exec(ctx, s, t)
	return err
}

func (p pgToken) Delete(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := "delete from tokens where id = $1;"
	_, err := p.pool.Exec(ctx, s, ID)
	return err
}

func (p pgToken) ListByUserID(ctx context.Context, UserID string) ([]models.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	s := "select data from tokens where user_id = $1;"
	rows, err := p.pool.Query(ctx, s, UserID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tokens := make([]models.Token, 0)
	for rows.Next() {
		var token models.Token
		err := rows.Scan(&token)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}
