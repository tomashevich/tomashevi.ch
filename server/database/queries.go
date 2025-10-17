package database

import "context"

func (d Database) GiveSoulToHel(ctx context.Context, seed, address string) error {
	if _, err := d.db.Exec("INSERT INTO souls (seed, address) VALUES (?, ?)", seed, address); err != nil {
		return err
	}

	return nil
}

func (d Database) GetSeed(ctx context.Context, id int) (string, error) {
	row := d.db.QueryRow("SELECT seed FROM souls WHERE id=?", id)
	var seed string

	if row.Err() != nil {
		return seed, row.Err()
	}

	if err := row.Scan(&seed); err != nil {
		return seed, err
	}

	return seed, nil
}

func (d Database) GetSoulIDByIP(ctx context.Context, address string) (int, error) {
	row := d.db.QueryRow("SELECT id FROM souls WHERE address=?", address)
	var id int

	if row.Err() != nil {
		return id, row.Err()
	}

	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (d Database) GetSeeds(ctx context.Context, limit, offset int64) ([]string, error) {
	rows, err := d.db.Query("SELECT seed FROM souls ORDER BY seed DESC LIMIT ? OFFSET ?", limit, offset)
	var seeds []string

	if err != nil {
		return seeds, err
	}

	defer rows.Close()

	for rows.Next() {
		var seed string
		if err := rows.Scan(&seed); err != nil {
			return nil, err
		}
		seeds = append(seeds, seed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seeds, nil
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

func (d Database) PaintPixel(ctx context.Context, soul_id, x, y int, color int) error {
	if _, err := d.db.Exec("INSERT INTO pixels (soul_id, color, x, y) VALUES (?, ?, ?, ?)", soul_id, color, x, y); err != nil {
		return err
	}

	return nil
}
