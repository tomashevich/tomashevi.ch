package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"

	"tomashevich/server"
	"tomashevich/server/database"
)

//go:embed static/*
//go:embed static/icons/*.svg
var staticFiles embed.FS

func main() {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("cant make fs from embeded fs %s", err.Error())
	}

	db, err := database.NewDatabase("storage.db")
	if err != nil {
		log.Fatalf("cant init storage with err %s", err.Error())
	}

	s := server.NewServer(":8037", staticFS, db)
	log.Fatal(s.Run())
}
