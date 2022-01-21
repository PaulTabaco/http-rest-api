package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"paulTabaco/http-rest-api/internal/app/model"
	"paulTabaco/http-rest-api/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName        = "myServerName"
	ctxKeyUser  CtxKey = iota
)

var (
	errIncorrectEmailOrPassord = errors.New("incorrect email or password")
	errNotAuthenticated        = errors.New("Not authenticated")
)

type CtxKey int8 // we use type here because it's good practice use special types istead simple (str or int)

type Server struct {
	router *mux.Router
	//logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *Server {
	s := &Server{
		router: mux.NewRouter(),
		//logger:       logrus.New(),
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
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
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
		//next.ServeHTTP(rw, r)
		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u))) // new context = current context + added values
	})
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
