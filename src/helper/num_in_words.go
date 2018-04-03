package helper

import (
	"errors"
	"math"
	"strings"
)

type number_string struct {
	name string
	rubl string
	kop  string
}

var number = map[int64]number_string{
	1:  {name: "один", rubl: "рубль", kop: "копейка"},
	2:  {name: "два", rubl: "рубля", kop: "копейки"},
	3:  {name: "три", rubl: "рубля", kop: "копейки"},
	4:  {name: "четыре", rubl: "рубля", kop: "копейки"},
	5:  {name: "пять", rubl: "рублей", kop: "копеек"},
	6:  {name: "шесть", rubl: "рублей", kop: "копеек"},
	7:  {name: "семь", rubl: "рублей", kop: "копеек"},
	8:  {name: "восемь", rubl: "рублей", kop: "копеек"},
	9:  {name: "девять", rubl: "рублей", kop: "копеек"},
	0:  {name: "", rubl: "рублей", kop: "копеек"},
	10: {name: "десять", rubl: "рублей", kop: "копеек"},
	11: {name: "одиндцать", rubl: "рублей", kop: "копеек"},
	12: {name: "двенадцать", rubl: "рублей", kop: "копеек"},
	13: {name: "тринадцать", rubl: "рублей", kop: "копеек"},
	14: {name: "четырнадцать", rubl: "рублей", kop: "копеек"},
	15: {name: "пятнадцать", rubl: "рублей", kop: "копеек"},
	16: {name: "шестнадцать", rubl: "рублей", kop: "копеек"},
	17: {name: "семнадцать", rubl: "рублей", kop: "копеек"},
	18: {name: "восемнадцать", rubl: "рублей", kop: "копеек"},
	19: {name: "девятнадцать", rubl: "рублей", kop: "копеек"},
}
var ten = map[int64]string{
	1: "десять",
	2: "двадцать",
	3: "тридцать",
	4: "сорок",
	5: "пятьдесят",
	6: "шестьдесят",
	7: "семьдесят",
	8: "восемьдесят",
	9: "девяносто",
}

var hundred = map[int64]string{
	1: "сто",
	2: "двести",
	3: "триста",
	4: "четыреста",
	5: "пятьсот",
	6: "шестьсот",
	7: "семьсот",
	8: "восемьсот",
	9: "девятьсот",
}

type Nc struct {
	Full []string
	Unit []string
	MF   []string
}

func Naming(code int) (Nc, error) {

	switch code {
	case -1:
		return Nc{
			Full: make([]string, 3),
			Unit: make([]string, 3),
			MF:   []string{"один", "два"},
		}, nil

	case -2:
		return Nc{
			Full: make([]string, 3),
			Unit: make([]string, 3),
			MF:   []string{"одна", "две"},
		}, nil

	case -3:
		return Nc{
			Full: make([]string, 3),
			Unit: make([]string, 3),
			MF:   []string{"одно", "два"},
		}, nil

	case -4:
		return Nc{
			Full: []string{"Минута", "Минуты", "Минут"},
			Unit: make([]string, 3),
			MF:   []string{"одна", "две"},
		}, nil

	}

	return Nc{}, errors.New("Не найден нейминг")

}

func NumInWords(fSum float64, currency int) (string, error) {

	naming, err := Naming(currency)

	if err != nil {

		return "", err

	}

	//отделяем рубли
	int, frac := math.Modf(fSum)

	iRubl := int64(int)

	fKop := frac * 100

	_, frac = math.Modf(fKop)

	if frac >= 0.5 {
		fKop = math.Ceil(fKop)
	} else {
		frac = math.Floor(fKop)
	}

	iKop := int64(fKop)

	//выходная строка
	sum := ""
	pos := 1
	//ограничение
	if iRubl > 999999999999999 {

		return sum, errors.New("Число должно быть меньше квадриллиона")

	}
	if iKop > 99 {

		return sum, errors.New("Получилось больше 99 копеек, что-то не так")

	}

	for iRubl > 0 {

		digit := iRubl % 10

		switch pos {
		//обрабатываем каждый первый разряд в тройке
		case 1, 4, 7, 10, 13:

			if iRubl%100/10 == 1 {

				digit = iRubl % 100

			}

			//тысячи
			if pos == 4 && iRubl%1000 > 0 {
				switch digit {
				case 1:
					sum = " тысяча" + sum
				case 2, 3, 4:
					sum = " тысячи" + sum
				default:
					sum = " тысяч" + sum
				}
			}
			//миллионы
			if pos == 7 {
				switch digit {
				case 1:
					sum = " миллион" + sum
				case 2, 3, 4:
					sum = " миллиона" + sum
				default:
					sum = " миллионов" + sum
				}
			}

			//миллиарды
			if pos == 10 {
				switch digit {
				case 1:
					sum = " миллиард" + sum
				case 2, 3, 4:
					sum = " миллиарда" + sum
				default:
					sum = " миллиардов" + sum
				}
			}

			//триллионы
			if pos == 13 {
				switch digit {
				case 1:
					sum = " триллион" + sum
				case 2, 3, 4:
					sum = " триллиона" + sum
				default:
					sum = " триллионов" + sum
				}
			}

			// добавляем название валюты в требуемой форме
			if sum == "" {

				switch iRubl {
				case 1:
					sum = " " + strings.ToLower(naming.Full[0])
				case 2, 3, 4:
					sum = " " + strings.ToLower(naming.Full[1])
				default:
					sum = " " + strings.ToLower(naming.Full[2])
				}

				//sum=" "+number[digit].rubl
			}

			//корректировка для тысяч
			if pos == 4 && digit == 1 {
				sum = " одна" + sum
			} else {
				if pos == 4 && digit == 2 {
					sum = " две" + sum
				} else {
					if digit != 0 {
						sum = " " + number[digit].name + sum
					}
				}
			}

			pos++
			if iRubl%100/10 == 1 {
				pos++
				iRubl = iRubl / 10
			}
		case 2, 5, 8, 11, 14:
			if digit != 0 {
				sum = " " + ten[digit] + sum
			}
			pos++
		case 3, 6, 9, 12, 15:
			if digit != 0 {
				sum = " " + hundred[digit] + sum
			}
			pos++
		}
		iRubl = iRubl / 10

	}

	sum = strings.Trim(sum, " ")

	//Для чисел в разных родах и времени
	if currency < 0 {

		namingM, _ := Naming(-1)

		if strings.HasSuffix(sum, namingM.MF[0]) {

			sum = strings.TrimSuffix(sum, namingM.MF[0]) + naming.MF[0]

		}

		if strings.HasSuffix(sum, namingM.MF[1]) {

			sum = strings.TrimSuffix(sum, namingM.MF[1]) + naming.MF[1]

		}

		//strings.Replace(sum, namingM[0], naming[0], )
	}

	switch {
	case iKop == 0:
		sum = sum
	case iKop > 3 && iKop <= 19:
		//sum=sum+ " "+number[iKop].name+" "+number[iKop].kop
		sum = sum + " " + number[iKop].name + " " + strings.ToLower(naming.Unit[2])
	case iKop%10 == 0:
		//sum=sum+ " "+ten[iKop/10]+" "+number[0].kop
		sum = sum + " " + ten[iKop/10] + " " + strings.ToLower(naming.Unit[2])
	case iKop%10 == 1:
		//sum=sum+ " "+ten[iKop/10]+" одна "+number[1].kop
		sum = sum + " " + ten[iKop/10] + " " + naming.MF[0] + " " + strings.ToLower(naming.Unit[1])
	case iKop%10 == 2:
		//sum=sum+ " "+ten[iKop/10]+" две "+number[2].kop
		sum = sum + " " + ten[iKop/10] + " " + naming.MF[1] + " " + strings.ToLower(naming.Unit[1])
	case iKop%10 == 3:
		//sum=sum+ " "+ten[iKop/10]+" две "+number[2].kop
		sum = sum + " " + ten[iKop/10] + " " + strings.TrimLeft(number[iKop].name+" "+strings.ToLower(naming.Unit[2]), " ")
	default:
		//sum=sum+" "+ten[iKop/10]+" "+number[iKop%10].name+" "+number[iKop%10].kop
		sum = sum + " " + ten[iKop/10] + " " + strings.TrimLeft(number[iKop%10].name+" "+strings.ToLower(naming.Unit[2]), " ")

	}

	return sum, nil

}
