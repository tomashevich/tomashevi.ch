package utils

import (
	"os"
)

type Config struct {
	DatabaseFile string            `json:"database_file"` // sqlite3 database file
	Server       ServerConfig      `json:"server"`
	RateLimiter  RateLimiterConfig `json:"ratelimiter"`
	Caches       CacheConfig       `json:"caches"`
}

type ServerConfig struct {
	Address      string `json:"address"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
}

type RateLimiterConfig struct {
	MaxRequests int `json:"max_requests"`
	InSeconds   int `json:"in_seconds"`
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
	if err := UnmarshalJSON(file, &config); err != nil {
		return config, err
	}

	return config, nil
}

func ParseConfigString(rawJson string) (Config, error) {
	var config Config
	if err := UnmarshalJSONString(rawJson, &config); err != nil {
		return config, err
	}

	return config, nil
}
