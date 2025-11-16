package v1

type Geo struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}
