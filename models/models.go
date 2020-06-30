package models

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

// version is the generation of the current models package
const modelVersion = 1

var now = time.Now
var entropy = ulid.Monotonic(rand.Reader, 0)

func NewID() string {
	return ulid.MustNew(ulid.Timestamp(now()), entropy).String()
}
