package apiserver

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIserver ...
type APIserver struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

// New ...
func New(config *Config) *APIserver {
	return &APIserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIserver) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("Starting API server")
	return nil
}

//Configure logger ...
func (s *APIserver) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.configureRouter()

	s.logger.SetLevel(level)

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

//Configure router ...
func (s *APIserver) configureRouter() {
	s.router.HandleFunc("/hellow", s.handleHello())
}

func (s *APIserver) handleHello() http.HandlerFunc {
	// Here can put some vars and this code executes ones only. Can set local types and avoide harm global

	// Start
	// ..
	// Here can put some vars and this code exutetes ones only. End

	// Other useal business logic
	return func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "HELLOW-OW-OW !!")
	}
}
