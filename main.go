package main

import (
	"embed"
	_ "embed"
	"log"

	"tomashevich/server"
)

//go:embed static/*
//go:embed static/icons/*.svg
var staticFiles embed.FS

func main() {
	server := server.NewServer(":8037", staticFiles)

	log.Fatal(server.Run())
}
