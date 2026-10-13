package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/health", s.handleHealth()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/users", s.handleUsersList()).Methods("GET")
	s.router.HandleFunc("/users/{id}", s.handleUserGet()).Methods("GET")
	s.router.HandleFunc("/users/{id}", s.handleUserUpdate()).Methods("PATCH")
	s.router.HandleFunc("/users/{id}", s.handleUserDelete()).Methods("DELETE")
}

// ---------- handlers ----------

func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   *int   `json:"age"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, store.ErrInvalidData)
			return
		}

		u := &model.User{
			Name:  req.Name,
			Email: req.Email,
			Age:   req.Age,
		}

		if err := s.store.User().Create(ctx, u); err != nil {
			switch {
			case errors.Is(err, store.ErrDataConflict):
				s.error(w, r, http.StatusConflict, err)
			case errors.Is(err, store.ErrInvalidData):
				s.error(w, r, http.StatusBadRequest, err)
			default:
				s.error(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleUsersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		limit, err := parseIntQuery(r, "limit", 20)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		offset, err := parseIntQuery(r, "offset", 0)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if limit <= 0 || limit > 100 {
			s.error(w, r, http.StatusBadRequest, errors.New("limit must be between 1 and 100"))
			return
		}
		if offset < 0 {
			s.error(w, r, http.StatusBadRequest, errors.New("offset must be >= 0"))
			return
		}

		users, err := s.store.User().GetAll(ctx, limit, offset)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) handleUserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := parseIDParam(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(w, r, http.StatusNotFound, err)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleUserUpdate() http.HandlerFunc {
	type request struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
		Age   *int    `json:"age"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := parseIDParam(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		upd := model.UserUpdate{
			Name:  req.Name,
			Email: req.Email,
			Age:   req.Age,
		}

		u, err := s.store.User().Update(ctx, id, upd)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				s.error(w, r, http.StatusNotFound, err)
			case errors.Is(err, store.ErrDataConflict):
				s.error(w, r, http.StatusConflict, err)
			default:
				s.error(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := parseIDParam(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().Delete(ctx, id); err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(w, r, http.StatusNotFound, err)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ---------- helpers ----------

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func parseIDParam(r *http.Request) (int, error) {
	raw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func parseIntQuery(r *http.Request, key string, def int) (int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.New("invalid " + key)
	}
	return v, nil
}
