package items

import (
	"time"
	"lib/auth"
)

type PTL struct {
	ItemData *Item
	C [][]timeData
}

type timeData struct{
	TypeID uint
	TypeName string
	ManufTime time.Duration
	Quantity uint
	Start time.Time
	End time.Time
}

func ProductionTimeline(itemId, regionId, batch uint, user *auth.User) (PTL, error) {

	res := PTL{}

	item, err := GetItem(itemId, regionId, user)
	if err != nil {

		return res, err

	}

	res.ItemData = item

	st := []timeData{}

	for _, i := range item.Materials {

		if i.ComponentBPO != 0 {

			subItem, err := GetItem(i.ComponentBPO, regionId, user)
			if err != nil {

				return res, err

			}

			st = append(st, timeData{
				TypeID: subItem.TypeID,
				TypeName: subItem.ProductName,
				ManufTime: time.Duration(float64(subItem.AdjTime) * i.Quantity),
				Quantity: uint(i.Quantity),
			})

		}

	}

	res.C = append(res.C, st)

	res.C = append(res.C, []timeData{{
		TypeID: item.TypeID,
		TypeName: item.ProductName,
		ManufTime: item.AdjTime,
		Quantity: batch,
	}})

	ptime := time.Now().UTC()
	maxDur := time.Nanosecond

	for i, r := range res.C {

		for j, r2 := range r {

			if maxDur < r2.ManufTime {
				maxDur = r2.ManufTime
			}

			res.C[i][j].Start = ptime
			res.C[i][j].End = ptime.Add(r2.ManufTime)
		}

		ptime = ptime.Add(maxDur)

	}

	return res, nil

}