package app

import (
	"github.com/flosch/pongo2"
	"helper"
	"app/crypto"
	"time"
	"fmt"
	"math"
)

func filters() {

	pongo2.RegisterFilter("NumInWords", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

		s, err := helper.NumInWords(in.Float(), param.Integer())

		if err != nil {

			s = err.Error()

		}

		return pongo2.AsValue(s), nil
	})

	pongo2.RegisterFilter("FormatFloat", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

		return pongo2.AsValue(helper.FormatNumber([]string{in.String()})), nil

	})

	pongo2.RegisterFilter("FormatPrice", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

		return pongo2.AsValue(helper.FormatPrice([]string{in.String()})), nil

	})

	pongo2.RegisterFilter("EncryptURL", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

		return pongo2.AsValue(crypto.EncryptUrl(uint(in.Integer()))), nil
	})

	pongo2.RegisterFilter("FormatDuration", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

		c := in.Float()

		d := time.Duration(int(c))

		if d.Hours() < 24 {

			return pongo2.AsValue(d.String()), nil

		}

		days := int(math.Floor(d.Hours() / 24))

		d = d - time.Duration(days * 86400) * time.Second

		return pongo2.AsValue(fmt.Sprintf("%dd%s", days ,d.String())), nil
	})

}
