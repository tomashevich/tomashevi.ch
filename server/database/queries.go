package database

func (d Database) AddFish(seed, address string) error {
	if _, err := d.db.Exec("INSERT INTO fishes (seed, address) VALUES (?, ?)", seed, address); err != nil {
		return err
	}

	return nil
}
