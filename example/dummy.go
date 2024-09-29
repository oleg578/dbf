package main

import (
	"encoding/json"
)

type Dummy struct {
	ID          int64       `json:"id"`
	Product     string      `json:"product"`
	Description interface{} `json:"description"`
	Price       float64     `json:"price"`
	Qty         int64       `json:"qty"`
	Date        string      `json:"date"`
}

func (d Dummy) Marshal() ([]byte, error) {
	return json.Marshal(d)
}
