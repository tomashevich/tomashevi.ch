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

func (d Database) GetFishes(ctx context.Context, limit, offset int64) ([]Fish, error) {
	rows, err := d.db.Query("SELECT * FROM fishes ORDER BY seed DESC LIMIT ? OFFSET ?", limit, offset)
	var fishes []Fish

	if err != nil {
		return fishes, err
	}

	defer rows.Close()

	for rows.Next() {
		var fish Fish
		if err := rows.Scan(&fish.Seed, &fish.Address); err != nil {
			return nil, err
		}
		fishes = append(fishes, fish)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fishes, nil
}

func (d Database) GetPixels(ctx context.Context) ([]Pixel, error) {
	rows, err := d.db.Query("SELECT * FROM pixels LIMIT 4000")
	var pixels []Pixel

	if err != nil {
		return pixels, err
	}

	defer rows.Close()

	for rows.Next() {
		var pixel Pixel
		if err := rows.Scan(&pixel.SoulId, &pixel.Color, &pixel.X, &pixel.Y); err != nil {
			return nil, err
		}
		pixels = append(pixels, pixel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pixels, nil
}

func (d Database) PaintPixel(ctx context.Context, soul_id, x, y int, color string) error {
	if _, err := d.db.Exec("INSERT INTO pixels (soul_id, color, x, y) VALUES (?, ?, ?, ?)", soul_id, color, x, y); err != nil {
		return err
	}

	return nil
}
