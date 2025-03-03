package server

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ed-henrique/tatu/internal/db"
	"github.com/ed-henrique/tatu/internal/models"
)

type Server struct {
	db *sql.DB
	q  *models.Queries

	Mux *http.ServeMux
}

type ServerOption func(*Server)

func WithDB(db *sql.DB) ServerOption {
	return func(s *Server) {
		s.db = db
	}
}

func New(options ...ServerOption) *Server {
	dbConn := db.New(":memory:")

	s := Server{
		db:  dbConn,
		q:   models.New(dbConn),
		Mux: http.NewServeMux(),
	}

	for _, o := range options {
		o(&s)
	}

	return &s
}

func (s *Server) Routes() {
	s.Mux.Handle("POST /secrets", s.addSecret())
}

func (s *Server) Run() {
	if err := http.ListenAndServe(":8080", s.Mux); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start server: %s", err.Error())
		os.Exit(1)
	}
}

func (s *Server) addSecret() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-type") != "application/json" {
			http.Error(w, "body not in json format", http.StatusBadRequest)
			return
		}

		dto := struct {
			Secret string `json:"secret,omitempty"`
		}{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if dto.Secret == "" {
			http.Error(w, "secret is empty", http.StatusBadRequest)
			return
		}

		convertedSecret, err := base64.URLEncoding.DecodeString(dto.Secret)
		if err != nil {
			http.Error(w, "secret is malformed", http.StatusBadRequest)
			return
		}

		id, err := s.q.AddSecret(r.Context(), convertedSecret)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "new secret: %d", id)
	})
}
