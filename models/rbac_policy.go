package models

import (
	"errors"
	"fmt"
	"time"
)

const (
	DefaultAdminRBAC  = Label("default-admin")
	DefaultGlobalRBAC = Label("default-global")
)

// RBACPolicy represents a policy defining role based access control
type RBACPolicy struct {
	ID        string    `json:"id"`
	Version   uint8     `json:"version"`
	Label     Label     `json:"label"`
	Spec      *RBACSpec `json:"spec"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate that the policy is valid
func (p RBACPolicy) Validate() error {
	if p.Version < 1 {
		return errors.New("version must be greater than zero")
	}

	if p.CreatedAt.IsZero() {
		return errors.New("createdAt cannot be a zero time")
	}

	if p.UpdatedAt.IsZero() {
		return errors.New("updatedAt cannot be a zero time")
	}

	if p.UpdatedAt.Before(p.CreatedAt) {
		return errors.New("updatedAt cannot be before CreatedAt")
	}

	err := p.Spec.Validate()
	if err != nil {
		return fmt.Errorf("policy has an invalid spec: %w", err)
	}

	return nil
}

// NewRBACPolicy returns new RBACPolicy
func NewRBACPolicy(label Label, spec *RBACSpec) RBACPolicy {
	return RBACPolicy{
		ID:        NewID(),
		Version:   modelVersion,
		Label:     label,
		Spec:      spec,
		CreatedAt: now(),
	}
}

// ParseRBACPolicy can convert a yaml document into a Policy
func ParseRBACPolicy(data []byte) (*RBACPolicy, error) {
	spec, err := ParseRBACSpec(data)
	if err != nil {
		return nil, err
	}

	rbac := NewRBACPolicy(spec.Label, spec)

	return &rbac, nil
}
