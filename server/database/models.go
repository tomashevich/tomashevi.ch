package database

import "time"

type Fish struct {
	Seed      string    `json:"seed"`
	Address   string    `json:"-"`
	SpawnedAt time.Time `json:"spawned_at"`
}
