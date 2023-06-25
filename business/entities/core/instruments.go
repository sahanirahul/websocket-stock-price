package core

type Instrument struct {
	Symbol         string  `json:"symbol"`
	Underlying     string  `json:"underlying"`
	Token          int64   `json:"token"`
	InstrumentType string  `json:"instrument_type"`
	Expiry         string  `json:"expiry"`
	Strike         int64   `json:"strike"`
	Price          float64 `json:"price"`
}
