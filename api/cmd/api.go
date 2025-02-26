package cmd

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ed-henrique/tatu/api/internal/db"
	"github.com/ed-henrique/tatu/api/internal/models"
)

type Server struct {
	db  *sql.DB
	q   *models.Queries
	mux *http.ServeMux
}

func Run() {
	dbConn := db.New("db.sqlite")
	s := Server{
		db:  dbConn,
		q:   models.New(dbConn),
		mux: http.NewServeMux(),
	}

	s.mux.Handle("POST /secrets", s.addSecret())

	if err := http.ListenAndServe(":8080", s.mux); err != nil {
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
