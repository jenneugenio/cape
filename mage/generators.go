package mage

import (
	"context"
)

type Generator interface {
	Generate(_ context.Context) error
}
