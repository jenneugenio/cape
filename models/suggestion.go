package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type SuggestionState int

func (ss *SuggestionState) UnmarshalGQL(v interface{}) error {
	switch s := v.(type) {
	case string:
		rawState, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		state := SuggestionState(rawState)
		*ss = state
		return nil
	default:
		return fmt.Errorf("error unmarshalling suggestion state")
	}
}

// MarshalGQL marshals a primitive.URL to a strong for GraphQL
func (ss SuggestionState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Itoa(int(ss)))
}

const (
	SuggestionPending SuggestionState = iota
	SuggestionApproved
	SuggestionRejected
)

type Suggestion struct {
	ID        string          `json:"id"`
	ProjectID string          `json:"project_id"`
	PolicyID  string          `json:"project_spec_id"`
	State     SuggestionState `json:"state"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
