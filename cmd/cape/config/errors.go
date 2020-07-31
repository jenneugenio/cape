package config

import (
	errors "github.com/capeprivacy/cape/partyerrors"
)

var (
	InvalidVersionCause     = errors.NewCause(errors.BadRequestCategory, "invalid_version")
	MissingConfigCause      = errors.NewCause(errors.BadRequestCategory, "missing_config")
	InvalidConfigCause      = errors.NewCause(errors.BadRequestCategory, "invalid_config")
	InvalidEnvCause         = errors.NewCause(errors.BadRequestCategory, "invalid_environment")
	InvalidPermissionsCause = errors.NewCause(errors.BadRequestCategory, "invalid_file_permissions")
	ExistingClusterCause    = errors.NewCause(errors.BadRequestCategory, "cluster_exists")
	ClusterNotFoundCause    = errors.NewCause(errors.NotFoundCategory, "cluster_not_found")
)
