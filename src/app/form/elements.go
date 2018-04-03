package form

import (
	"fmt"
	"strconv"
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
	Options Options
}

func (e *Element) GetValue() string {

	val := fmt.Sprintf("%v", e.Value)

	switch e.Value.(type) {

	case float32, float64:

		val2, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return val
		}

		return fmt.Sprintf("%.2f", val2)

	case int, int8, int16, int32, int64:

		val2, err := strconv.ParseInt(val, 10,64)
		if err != nil {
			return val
		}

		return fmt.Sprintf("%d", val2)

	case  uint, uint8, uint16, uint32, uint64:

		val2, err := strconv.ParseUint(val, 10,64)
		if err != nil {
			return val
		}

		return fmt.Sprintf("%d", val2)

	default:

		return val

	}


}

func (e *Element) TplType() string {
	return e.TypeName
}