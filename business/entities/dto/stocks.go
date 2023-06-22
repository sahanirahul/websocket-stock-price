package dto

type UnderLyings struct {
	Symbol         string `json:"symbol"`
	Underlying     string `json:"underlying"`
	Token          string `json:"token"`
	InstrumentType string `json:"instrument_type"`
	Expiry         string `json:"expiry"`
	Strike         string `json:"strike"`
}
