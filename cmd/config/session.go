package config

// Session stores CLI "session" data that overrides data in the config.
// For example, when a command has the flag --cluster it overrides the
// value in the config.
type Session struct {
	config  *Config
	cluster *Cluster
}

// NewSession returns a new session
func NewSession(config *Config, cluster *Cluster) *Session {
	return &Session{
		config:  config,
		cluster: cluster,
	}
}

// Cluster returns the cluster stored on session if exists or the
// config cluster if it doesn't exist.
func (s *Session) Cluster() (*Cluster, error) {
	if s.cluster == nil {
		cluster, err := s.config.Cluster()
		if err != nil {
			return nil, err
		}

		return cluster, nil
	}

	return s.cluster, nil
}
