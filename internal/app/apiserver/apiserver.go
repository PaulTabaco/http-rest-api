package apiserver

// APIserver ...
type APIserver struct {
	config *Config
}

// New ...
func New(config *Config) *APIserver {
	return &APIserver{}
}

// Start ...
func (s *APIserver) Start() error {
	return nil
}
