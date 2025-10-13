package database

type Fish struct {
	Seed    string `json:"seed"`
	Address string `json:"-"` // TODO: SOUL REFERENCE
}

type Soul struct {
	Id            int    `json:"id"`
	Address       string `json:"-"`
	PaintedPixels int    `json:"painted_pixels"`
}

type Pixel struct {
	SoulId int    `json:"soul_id"`
	Color  string `json:"color"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}
