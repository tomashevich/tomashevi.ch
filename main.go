package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"

	"tomashevich/server"
	"tomashevich/server/database"
	"tomashevich/utils"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed config.json
var rawConfig string

func main() {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("cant make fs from embeded fs %s", err.Error())
	}

	config, err := utils.ParseConfigString(rawConfig)
	if err != nil {
		log.Fatalf("cant load config file with err %s", err)
	}

	db, err := database.NewDatabase(config.DatabaseFile)
	if err != nil {
		log.Fatalf("cant init storage with err %s", err.Error())
	}

	s := server.NewServer(&config, staticFS, db)
	log.Fatal(s.Run())
}
