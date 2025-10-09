package server

import (
	"io/fs"
	"log"
	"net/http"

	"tomashevich/server/database"
	"tomashevich/server/handlers/api"
)

type Server struct {
	addr        string
	database    *database.Database
	staticFiles fs.FS
}

func NewServer(addr string, staticFiles fs.FS, database *database.Database) *Server {
	return &Server{
		addr,
		database,
		staticFiles,
	}
}

func (s Server) Run() error {
	router := http.NewServeMux()

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	// Register static files
	router.Handle("/", http.FileServerFS(s.staticFiles))

	// Register API handlers
	api.RegisterFishes(router)

	log.Printf("starting server at %s", s.addr)

	return server.ListenAndServe()
}
