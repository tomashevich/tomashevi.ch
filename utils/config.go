package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address      string      `json:"address"`       // server addr
	DatabaseFile string      `json:"database_file"` // sqlite3 database file
	Caches       CacheConfig `json:"caches"`
}

type CacheConfig struct {
	StaticFiles int `json:"static_files"` // cache for static files
	FishesMe    int `json:"fish_me"`      // cache for GET fish:me
	PixelsLimit int `json:"pixels_limit"` // cache for limit of pixels
}

func ParseConfig(fileName string) (Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func ParseConfigString(rawJson string) (Config, error) {
	var config Config

	if err := json.Unmarshal([]byte(rawJson), &config); err != nil {
		return config, err
	}

	return config, nil
}
