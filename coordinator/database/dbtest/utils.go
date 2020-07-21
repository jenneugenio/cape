package dbtest

import (
	"crypto/rand"

	base32 "github.com/manifoldco/go-base32"
)

var dbNameByteLength = 4

func GenerateName() (string, error) {
	value := make([]byte, dbNameByteLength)
	_, err := rand.Read(value)
	if err != nil {
		return "", err
	}

	return "testdb_" + base32.EncodeToString(value), nil
}
