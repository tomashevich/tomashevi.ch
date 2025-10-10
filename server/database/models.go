package database

type Fish struct {
	Seed    string `json:"seed"`
	Address string `json:"-"`
}
