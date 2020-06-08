package database

import (
	"testing"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/types"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestPrimitive(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		name  string
		fn    func() (*Primitive, error)
		cause *errors.Cause
	}{
		{
			name: "valid primitive",
			fn: func() (*Primitive, error) {
				return NewPrimitive(types.TestMutable)
			},
		},
		{
			name: "invalid id value",
			fn: func() (*Primitive, error) {
				p, err := NewPrimitive(types.TestMutable)
				if err != nil {
					return nil, err
				}

				p.ID = EmptyID
				return p, nil
			},
			cause: &InvalidIDCause,
		},
		{
			name: "invalid version",
			fn: func() (*Primitive, error) {
				p, err := NewPrimitive(types.TestMutable)
				if err != nil {
					return nil, err
				}

				p.Version = 0
				return p, nil
			},
			cause: &InvalidVersionCause,
		},
		{
			name: "invalid createdat",
			fn: func() (*Primitive, error) {
				p, err := NewPrimitive(types.TestMutable)
				if err != nil {
					return nil, err
				}

				p.CreatedAt = time.Time{}
				return p, nil
			},
			cause: &InvalidTimeCause,
		},
		{
			name: "invalid updated at",
			fn: func() (*Primitive, error) {
				p, err := NewPrimitive(types.TestMutable)
				if err != nil {
					return nil, err
				}

				p.UpdatedAt = time.Time{}
				return p, nil
			},
			cause: &InvalidTimeCause,
		},
		{
			name: "updated at before created at",
			fn: func() (*Primitive, error) {
				p, err := NewPrimitive(types.TestMutable)
				if err != nil {
					return nil, err
				}

				p.UpdatedAt = p.CreatedAt.Add(-1 * time.Minute)
				return p, nil
			},
			cause: &InvalidTimeCause,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := tc.fn()
			gm.Expect(err).To(gm.BeNil())

			err = p.Validate()
			if tc.cause != nil {
				gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(err).To(gm.BeNil())
		})
	}
}
