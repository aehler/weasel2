package form

import (
	"fmt"
)

type Element struct {
	Name string
	HashName string `json:"n"`
	Value interface {} `json:"v"`
	Label string `json:"l"`
	PlaceHolder string
	Description string
	Type uint `json:"-"`
	TypeName string `json:"t"`
	Order uint
	Options Options `json:"o"`
}

func (e *Element) GetValue() string {
	
	return fmt.Sprintf("%v", e.Value)
	
}
