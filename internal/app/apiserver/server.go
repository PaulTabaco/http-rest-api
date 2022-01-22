package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "myServerName"
	ctxKeyUser  CtxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassord = errors.New("incorrect email or password")
	errNotAuthenticated        = errors.New("not authenticated")
)

type CtxKey int8 // we use type here because it's good practice use special types istead simple (str or int)

type Server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *Server {
	s := &Server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}
	s.configureRouter()

	return s
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

func (s *Server) configureRouter() {
	s.router.Use(s.setRequestID) // Middlware for set ID to each request
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}))) // Middlware adds origins for relax* browsers (for all resourses) - (Access-Control-Allow-Origin: *)

	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	privateRouter := s.router.PathPrefix("/private").Subrouter() //subrouter for space for outhenticated only
	privateRouter.Use(s.authenticateUser)                        // Middlware check user by coockie
	privateRouter.HandleFunc("/whoami", s.handleWhoami()).Methods("Get")
}

func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		rw.Header().Set("X-Request-ID", id)
		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{ //made own local logger with special fields
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		// Format - started time, Method, endpoint
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		startTime := time.Now()
		rwWithCode := &ResponseWriter{rw, http.StatusOK}
		next.ServeHTTP(rwWithCode, r)
		logger.Infof("completed with %d %s in %v", rwWithCode.code, http.StatusText(rwWithCode.code), time.Since(startTime))
	})
}

func (s *Server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(rw, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		u, err := s.store.User().FindById(id.(int))
		if err != nil {
			s.error(rw, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u))) // new context = current context + added values
	})
}

func (s *Server) handleWhoami() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		s.respond(rw, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *Server) handleUsersCreate() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &Request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(rw, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		s.respond(rw, r, http.StatusCreated, u)
	}
}

func (s *Server) handleSessionsCreate() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		req := &Request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(rw, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			s.error(rw, r, http.StatusUnauthorized, errIncorrectEmailOrPassord)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, rw, session); err != nil {
			s.error(rw, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(rw, r, http.StatusOK, nil)
	}
}

func (s *Server) error(rw http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(rw, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	rw.WriteHeader(code)
	if data != nil {
		json.NewEncoder(rw).Encode(data)
	}
}
