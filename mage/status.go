package mage

import (
	"fmt"
)

type Status string

func (s Status) String() string {
	return string(s)
}

var (
	Running    Status = "running"
	Created    Status = "created"
	Restarting Status = "restarting"
	Removing   Status = "removing"
	Paused     Status = "paused"
	Exited     Status = "exited"
	Dead       Status = "dead"
	Unknown    Status = "unknown"
)

func ToStatus(in string) (Status, error) {
	switch in {
	case "running":
		return Running, nil
	case "created":
		return Created, nil
	case "Restarting":
		return Restarting, nil
	case "Removing":
		return Removing, nil
	case "Paused":
		return Paused, nil
	case "Exited":
		return Exited, nil
	case "Dead":
		return Dead, nil
	default:
		return Dead, fmt.Errorf("Unknown status: %s", in)
	}
}
