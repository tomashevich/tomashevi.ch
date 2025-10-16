package server

import (
	"io/fs"
	"log"
	"net/http"

	"tomashevich/server/database"
	"tomashevich/server/handler"
	"tomashevich/server/middleware"
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

	stack := middleware.MiddlewareStack(
		middleware.Helheim(s.database),
		middleware.Gzip(),
	)

	server := http.Server{
		Addr:    s.addr,
		Handler: stack(router),
	}

	// Register static files
	router.Handle("/", http.FileServerFS(s.staticFiles))

	// Register API handler
	handler.RegisterFishes(router, s.database)
	handler.RegisterPixels(router, s.database)

	log.Printf("starting server at %s", s.addr)

	return server.ListenAndServe()
}
