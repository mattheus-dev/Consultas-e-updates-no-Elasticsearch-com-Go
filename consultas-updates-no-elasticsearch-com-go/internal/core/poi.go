package core

type POI struct {
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
}

type POIRepository interface {
	GetByID(id string) (*POI, error)
	SearchByField(field, value string) ([]POI, error)
	Insert(poi *POI) (string, error)
	Update(id string, poi *POI) error
}
