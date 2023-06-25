package core

import "encoding/json"

type Instrument struct {
	Symbol         string  `json:"symbol"`
	Underlying     string  `json:"underlying"`
	Token          int64   `json:"token"`
	InstrumentType string  `json:"instrument_type"`
	Expiry         string  `json:"expiry"`
	Strike         int64   `json:"strike"`
	Price          float64 `json:"price"`
}

func (ins Instrument) MarshalBinary() (data []byte, err error) {
	return json.Marshal(ins)
}

func (ins *Instrument) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &ins)
}
