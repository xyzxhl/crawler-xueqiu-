package pj

import (
	"encoding/json"
	"time"
)

type LData struct {
	Column []string    `json:"column"`
	Item   [][]float32 `json:"item"`
	Symbol string      `json:"symbol"`
}

type LJson struct {
	Data    LData  `json:"data"`
	ErrCode int    `json:"error_code"`
	ErrDesc string `json:"error_description"`
}

type LChange struct {
	T       time.Time
	Percent float32
	Close   float32
}

type RFund struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type RData struct {
	Count int     `json:"count"`
	Funds []RFund `json:"list"`
}

type RJson struct {
	Data    RData  `json:"data"`
	ErrCode int    `json:"error_code"`
	ErrDesc string `json:"error_description"`
}

type UnknownError struct {
	errDesc string
}

func (ed UnknownError) Error() string {
	return ed.errDesc
}

func LProcessJson(content []byte) ([]LChange, error) {
	var tmp LJson
	if err := json.Unmarshal(content, &tmp); err != nil {
		return nil, err
	}

	if tmp.ErrCode != 0 {
		return nil, UnknownError{errDesc: tmp.ErrDesc}
	}

	data := tmp.Data.Item[:]
	change := make([]LChange, len(data))
	for i, v := range data {
		change[i].T = time.Unix(int64(v[0]/1000)+86400, 0)
		change[i].Percent = v[7]
		change[i].Close = v[5]
	}

	return change, nil
}

func RProcessJson(content []byte) ([]RFund, error) {
	var tmp RJson
	if err := json.Unmarshal(content, &tmp); err != nil {
		return nil, err
	}

	if tmp.ErrCode != 0 {
		return nil, UnknownError{errDesc: tmp.ErrDesc}
	}

	data := tmp.Data.Funds[:]
	return data, nil
}
