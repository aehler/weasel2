package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func FormatNumber(s []string) string {

	f, err := strconv.ParseFloat(strings.Join(s, ""), 64)

	if err != nil {

		return strings.Join(s, " ")

	}

	b := []rune(fmt.Sprintf("%.6f", f))

	nb := []rune{}

	delimiterIndex := strings.Index(string(b), ",")
	if delimiterIndex == -1 {
		delimiterIndex = strings.Index(string(b), ".")
	}

	fb := b[:delimiterIndex]

	for i := len(fb) - 1; i >= 0; i-- {

		nb = append(nb, fb[i])

		if (len(fb)-i)%3 == 0 {

			nb = append(nb, ' ')

		}

	}

	b2 := []rune{}

	for i := len(nb); i > 0; i-- {

		b2 = append(b2, nb[i-1])

	}

	b2 = append(b2, b[delimiterIndex:]...)

	st := strings.TrimRight(string(b2), "0")

	return strings.TrimRight(strings.Replace(strings.Trim(st, " "), ".", ",", -1), ",")
}

func FormatPrice(s []string) string {

	f, err := strconv.ParseFloat(strings.Join(s, ""), 64)

	if err != nil {

		return strings.Join(s, " ")

	}

	b := []rune(fmt.Sprintf("%.2f", f))

	nb := []rune{}

	for i := len(b); i > 0; i-- {

		nb = append(nb, b[i-1])

		if (len(b)-i)%3 == 2 && len(b)-i > 4 {

			nb = append(nb, ' ')

		}

	}

	b = []rune{}

	for i := len(nb); i > 0; i-- {

		b = append(b, nb[i-1])

	}

	//Trimming prcision when no prcision
	if string(b[len(b)-3:]) == ".00" {

		b = b[:len(b)-3]

	}

	return strings.TrimRight(strings.Replace(strings.Trim(string(b), " "), ".", ",", -1), ",")
}
