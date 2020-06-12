// +build ui

package ui

import (
	"net/http"

	"github.com/markbates/pkger"
)

var Version = "unknown"

func Handler() http.Handler {
	fs := pkger.Dir("/coordinator/ui/assets/dist")
	return http.FileServer(fs)
}
