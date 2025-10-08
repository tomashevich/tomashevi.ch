package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"tomashevich/server/database"
	"tomashevich/server/handlers/api"
)

type Server struct {
	addr        string
	database    *database.Database
	staticFiles embed.FS
}

func NewServer(addr string, staticFiles embed.FS, database *database.Database) *Server {
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

	staticFS, err := fs.Sub(s.staticFiles, "static")
	if err != nil {
		return err
	}

	// Register static files
	router.Handle("/", http.FileServerFS(staticFS))

	// Register API handlers
	api.RegisterFishes(router)

	log.Printf("starting server at %s", s.addr)

	return server.ListenAndServe()
}
