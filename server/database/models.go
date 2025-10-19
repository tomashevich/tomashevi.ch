package database

type Soul struct {
	Id            int    `json:"id"`
	Address       string `json:"-"`
	Seed          string `json:"seed"`
	PaintedPixels int    `json:"painted_pixels"`
}

type Pixel struct {
	SoulId int `json:"soul_id"`
	Color  int `json:"color"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

// For init field query
type PixelPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}
