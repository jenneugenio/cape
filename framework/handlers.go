package framework

import (
	"net/http"

	"github.com/capeprivacy/cape/version"
)

// VersionResponse represents the data returned when querying the version
// handler
type VersionResponse struct {
	InstanceID string `json:"instance_id"`
	Version    string `json:"version"`
	BuildDate  string `json:"build_date"`
}

// VersionHandler returns the version information for this instance of cape.
func VersionHandler(instanceID string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, &VersionResponse{
			InstanceID: instanceID,
			Version:    version.Version,
			BuildDate:  version.BuildDate,
		})
	})
}
