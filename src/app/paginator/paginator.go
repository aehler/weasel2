package paginator

import (
	"math"
	"net/url"
	"strconv"
)

const (
	size uint = 5
)

func NewPaginator(current, total, limit uint, path string, params map[string]string) Paginator {

	totalPages := uint(math.Ceil(float64(total) / float64(limit)))

	var (
		start    uint = 1
		stop     uint = totalPages
		rawQuery string
	)

	if params != nil {

		values := url.Values{}

		for key, value := range params {

			values.Set(key, value)
		}

		rawQuery = values.Encode()
	}

	if totalPages > size*2 {

		if current > size {

			start = current - size
		}

		stop = start + size*2 - 1

		if stop > totalPages {

			stop = totalPages
		}

		if stop-start < size*2 {

			start = stop - size*2 + 1
		}
	}

	list := make([]Page, stop-start+1)

	z := 0

	for i := start; i <= stop; i++ {

		list[z].Current = i == current
		list[z].Number = i

		if i == 1 {

			if rawQuery != "" {

				list[z].URL = path + "?" + rawQuery

			} else {

				list[z].URL = path
			}

		} else {

			if rawQuery != "" {

				list[z].URL = path + "page/" + strconv.Itoa(int(i)) + "/?" + rawQuery

			} else {

				list[z].URL = path + "page/" + strconv.Itoa(int(i)) + "/"
			}
		}

		z++
	}

	p := Paginator{
		TotalPages: totalPages,
		TotalItems: total,
		Current:    current,
		List:       list,
	}

	if p.TotalPages > 1 && current > 1 {

		if rawQuery != "" {

			p.Prev = path + "page/" + strconv.Itoa(int(current-1)) + "/?" + rawQuery

			p.First = path + "?" + rawQuery

		} else {

			p.Prev = path + "page/" + strconv.Itoa(int(current-1))

			p.First = path
		}
	}

	if p.TotalPages > 1 && current < p.TotalPages {

		if rawQuery != "" {

			p.Next = path + "page/" + strconv.Itoa(int(current+1)) + "/?" + rawQuery

			p.Last = path + "page/" + strconv.Itoa(int(p.TotalPages)) + "/?" + rawQuery

		} else {

			p.Next = path + "page/" + strconv.Itoa(int(current+1))

			p.Last = path + "page/" + strconv.Itoa(int(p.TotalPages))
		}
	}

	return p
}

type Page struct {
	Number  uint
	URL     string
	Current bool
}

type Paginator struct {
	Current    uint
	Prev       string
	Next       string
	TotalPages uint
	TotalItems uint
	List       []Page
	Last       string
	First      string
}
