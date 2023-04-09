package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aheld/listservice/db"
	"github.com/aheld/listservice/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	database, err := db.Initialize(os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatalf("Could not connect to DB %v", err)
	}
	defer database.Conn.Close()

	s := CreateNewServer(database)
	// mount probes
	s.MountInfrastructureHandlers()
	// mount api
	s.MountHandlers()
	http.ListenAndServe("localhost:3333", s.Router)
}

func CreateNewServer(database domain.ListDb) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.Database = database
	return s
}

func (s *Server) MountInfrastructureHandlers() {
	r := s.Router

	// Don't start taking traffic until the DB is migrated to the right version

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Heartbeat("/livez"))
	r.Use(middleware.Heartbeat("/readyz"))

	r.Get("/startz", s.CheckDbVersion)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
}

func (s *Server) CheckDbVersion(w http.ResponseWriter, r *http.Request) {
	msg, err := s.Database.CheckVersion()
	if err != nil {
		render.Render(w, r, ErrSrvUnavailable(err))
		return
	}
	w.Write([]byte(msg))
}

type Server struct {
	Router   *chi.Mux
	Database domain.ListDb
}
