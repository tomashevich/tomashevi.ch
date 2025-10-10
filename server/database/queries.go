package database

import "context"

func (d Database) AddFish(ctx context.Context, seed, address string) error {
	if _, err := d.db.Exec("INSERT INTO fishes (seed, address) VALUES (?, ?)", seed, address); err != nil {
		return err
	}

	return nil
}

func (d Database) GetFishByIP(ctx context.Context, address string) (Fish, error) {
	row := d.db.QueryRow("SELECT * FROM fishes WHERE address=?", address)
	var fish Fish

	if row.Err() != nil {
		return fish, row.Err()
	}

	if err := row.Scan(&fish.Seed, &fish.Address); err != nil {
		return fish, err
	}

	return fish, nil
}
