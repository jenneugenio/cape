package auth

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	b64 "github.com/manifoldco/go-base64"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantVal *b64.Value
		wantErr error
	}{
		{
			name:    "valid header returns correct token",
			in:      fmt.Sprintf("Bearer %s", base64.RawURLEncoding.EncodeToString([]byte("token"))),
			wantVal: b64.New([]byte("token")),
			wantErr: nil,
		},
		{
			name:    "header with too few parts errors",
			in:      "Bearer",
			wantVal: nil,
			wantErr: ErrorInvalidAuthHeader,
		},
		{
			name:    "header with too many parts errors",
			in:      "Bearer with too many parts",
			wantVal: nil,
			wantErr: ErrorInvalidAuthHeader,
		},
		{
			name:    "fails if token parts is not base64 encoded value",
			in:      "Bearer bad_token",
			wantVal: nil,
			wantErr: ErrorInvalidAuthHeader,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotVal, gotErr := GetBearerToken(test.in)
			if !reflect.DeepEqual(gotVal, test.wantVal) {
				t.Errorf("unexpected value returned: want %v got %v", test.wantVal, gotVal)
			}
			if gotErr != test.wantErr {
				t.Errorf("unexpected error returned: want %v got %v", test.wantErr, gotErr)
			}
		})
	}
}
