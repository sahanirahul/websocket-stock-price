package core

import (
	"encoding/json"
	"sensibull/stocks-api/business/entities/dto"

	"github.com/emirpasic/gods/sets/hashset"
)

type Instrument struct {
	Symbol         string  `json:"symbol"`
	Underlying     string  `json:"underlying"`
	Token          int64   `json:"token"`
	InstrumentType string  `json:"instrument_type"`
	Expiry         string  `json:"expiry"`
	Strike         float64 `json:"strike"`
	Price          float64 `json:"price"`
}

func (ins Instrument) MarshalBinary() (data []byte, err error) {
	return json.Marshal(ins)
}

func (ins *Instrument) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &ins)
}

func GetDtoInstrument(ci Instrument) dto.Instrument {
	return dto.Instrument{
		Symbol:         ci.Symbol,
		Token:          ci.Token,
		Underlying:     ci.Underlying,
		Expiry:         ci.Expiry,
		Strike:         ci.Strike,
		Price:          ci.Price,
		InstrumentType: ci.InstrumentType,
	}
}

func GetDtoInstruments(coreInstruments []Instrument) []dto.Instrument {
	dtos := []dto.Instrument{}
	for _, val := range coreInstruments {
		dtos = append(dtos, GetDtoInstrument(val))
	}
	return dtos
}

type Tokens struct {
	Set *hashset.Set
}

func NewTokenSet() Tokens {
	return Tokens{Set: hashset.New()}
}

func (tks *Tokens) MarshalBinary() (data []byte, err error) {
	return json.Marshal(tks)
}

func (tks *Tokens) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tks)
}

type WebsocketSubscription struct {
	MessageCommand string  `json:"msg_command"`
	DataType       string  `json:"data_type"`
	Tokens         []int64 `json:"tokens"`
}

type WebsocketPriceEvent struct {
	DataType string `json:"data_type"`
	Payload  struct {
		Token int64   `json:"token"`
		Price float64 `json:"price"`
	} `json:"payload"`
}
