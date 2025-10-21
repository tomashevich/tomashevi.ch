package server

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"tomashevich/server/database"
	"tomashevich/server/handler"
	"tomashevich/server/middleware"
	"tomashevich/utils"
)

type Server struct {
	config      *utils.Config
	database    *database.Database
	staticFiles fs.FS
}

func NewServer(config *utils.Config, staticFiles fs.FS, database *database.Database) *Server {
	return &Server{
		config,
		database,
		staticFiles,
	}
}

func (s Server) Run() error {
	router := http.NewServeMux()

	stack := middleware.MiddlewareStack(
		middleware.Helheim(s.database),
		middleware.Compress(),
		middleware.NewRateLimiter(s.config.RateLimiter.MaxRequests, time.Duration(s.config.RateLimiter.InSeconds)*time.Second).Middleware(),
	)

	server := http.Server{
		Addr:         s.config.Server.Address,
		Handler:      stack(router),
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}

	// Register static files
	router.Handle("/", middleware.Cache(time.Second*time.Duration(s.config.Caches.StaticFiles))(http.FileServerFS(s.staticFiles)))

	// Register API handler
	handler.RegisterFishes(router, s.database, &s.config.Caches)
	handler.RegisterPixels(router, s.database, &s.config.Caches)

	log.Printf("starting server at %s", s.config.Server.Address)

	return server.ListenAndServe()
}
