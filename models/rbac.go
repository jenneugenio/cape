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

// RBAC is a single defined policy
type RBAC struct {
	ID        string    `json:"id"`
	Version   uint8     `json:"version"`
	Label     Label     `json:"label"`
	Spec      *RBACSpec `json:"spec"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate that the policy is valid
func (p RBAC) Validate() error {
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

// NewRBAC returns a mutable policy struct
func NewRBAC(label Label, spec *RBACSpec) RBAC {
	return RBAC{
		ID:        NewID(),
		Version:   modelVersion,
		Label:     label,
		Spec:      spec,
		CreatedAt: now(),
	}
}

// ParseRBAC can convert a yaml document into a Policy
func ParseRBAC(data []byte) (*RBAC, error) {
	spec, err := ParseRBACSpec(data)
	if err != nil {
		return nil, err
	}

	rbac := NewRBAC(spec.Label, spec)

	return &rbac, nil
}
