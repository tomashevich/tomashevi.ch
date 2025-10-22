package database

import (
	"context"
	"strings"
)

func (d Database) GiveSoulToHel(ctx context.Context, seed, address string) (int, error) {
	row := d.db.QueryRowContext(ctx, "INSERT INTO souls (seed, address) VALUES (?, ?) RETURNING id", seed, address)

	var id int
	if row.Err() != nil {
		return id, row.Err()
	}

	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil

}

func (d Database) GetSeed(ctx context.Context, id int) (string, error) {
	row := d.db.QueryRowContext(ctx, "SELECT seed FROM souls WHERE id=?", id)
	var seed string

	if row.Err() != nil {
		return seed, row.Err()
	}

	if err := row.Scan(&seed); err != nil {
		return seed, err
	}

	return seed, nil
}

func (d Database) GetSoul(ctx context.Context, id int) (Soul, error) {
	row := d.db.QueryRowContext(ctx, "SELECT * FROM souls WHERE id=?", id)

	var soul Soul
	if row.Err() != nil {
		return soul, row.Err()
	}

	if err := row.Scan(&soul.Id, &soul.Address, &soul.Seed, &soul.PaintedPixels); err != nil {
		return soul, err
	}

	return soul, nil
}

func (d Database) GetSoulIDByIP(ctx context.Context, address string) (int, error) {
	row := d.db.QueryRowContext(ctx, "SELECT id FROM souls WHERE address=?", address)

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
	rows, err := d.db.QueryContext(ctx, "SELECT * FROM pixels")
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
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "UPDATE souls SET painted_pixels = painted_pixels + 1 WHERE id=?", soul_id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE pixels SET soul_id=?, color=? WHERE x=? AND y=?", soul_id, color, x, y); err != nil {
		return err
	}

	return tx.Commit()
}

func (d Database) IsPixelFieldInited(ctx context.Context) (bool, error) {
	row := d.db.QueryRowContext(ctx, "SELECT x FROM pixels LIMIT 1")

	var x int
	if err := row.Scan(&x); err != nil {
		return false, err
	}

	return true, nil
}

func (d Database) InitPixelField(ctx context.Context, pixels []PixelPosition, soulID, color int) error {
	if len(pixels) == 0 {
		return nil
	}

	// we dont use static soulID and color bc zt
	valuePlaceholder := "(?, ?, ?, ?)"
	placeholders := make([]string, 0, len(pixels))

	for range pixels {
		placeholders = append(placeholders, valuePlaceholder)
	}

	args := make([]any, 0, len(pixels)*4)
	for _, p := range pixels {
		args = append(args, soulID, color, p.X, p.Y)
	}

	baseSQL := "INSERT OR REPLACE INTO pixels (soul_id, color, x, y) VALUES "
	fullSQL := baseSQL + strings.Join(placeholders, ", ")

	_, err := d.db.ExecContext(ctx, fullSQL, args...)
	return err
}
