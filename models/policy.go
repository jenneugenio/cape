package models

import (
	"errors"
	"fmt"
	"time"
)

const (
	DefaultAdminPolicy  = Label("default-admin")
	DefaultGlobalPolicy = Label("default-global")
)

// Policy is a single defined policy
type Policy struct {
	ID        string      `json:"id"`
	Version   uint8       `json:"version"`
	Label     Label       `json:"label"`
	Spec      *PolicySpec `json:"spec"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Validate that the policy is valid
func (p Policy) Validate() error {
	if p.Version < 1 {
		return errors.New("Version must be greater than zero")
	}

	if p.CreatedAt.IsZero() {
		return errors.New("CreatedAt cannot be a zero time")
	}

	if p.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt cannot be a zero time")
	}

	if p.UpdatedAt.Before(p.CreatedAt) {
		return errors.New("UpdatedAt cannot be before CreatedAt")
	}

	err := p.Spec.Validate()
	if err != nil {
		return fmt.Errorf("policy has an invalid spec: %w", err)
	}

	return nil
}

// NewPolicy returns a mutable policy struct
func NewPolicy(label Label, spec *PolicySpec) Policy {
	return Policy{
		ID:        NewID(),
		Version:   modelVersion,
		Label:     label,
		Spec:      spec,
		CreatedAt: now(),
	}
}

// ParsePolicy can convert a yaml document into a Policy
func ParsePolicy(data []byte) (*Policy, error) {
	ps, err := ParsePolicySpec(data)
	if err != nil {
		return nil, err
	}

	policy := &Policy{
		Label: ps.Label,
		Spec:  ps,
	}

	return policy, nil
}
