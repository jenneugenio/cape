package mage

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver"
)

type Version struct {
	version   *semver.Version
	buildDate time.Time
}

func NewVersion(in string) (*Version, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return nil, fmt.Errorf("Could not parse version, invalid semver: %s", err.Error())
	}

	return &Version{
		version:   v,
		buildDate: time.Now().UTC(),
	}, nil
}

func (v *Version) Version() string {
	return v.version.String()
}

func (v *Version) BuildDate() string {
	return v.buildDate.Format(time.UnixDate)
}
