package market

import (
	"lib/esi/response"
	"lib/esi"
	"app/registry"
	"time"
	"fmt"
)

type MarketHistoryRow struct {
	RegionID uint `db:"regionID"`
	TypeID uint `db:"typeID"`
	response.MHItem
}

type Avgr struct {
	Avg float64
	Buy float64
	Sell float64
}

func GetRegionalAverage(itemID, regionID uint, daysAvg int) (res Avgr, err error) {

	ex := false
	if err := registry.Registry.Connect.SQLX().Get(&ex, `select exists (select 1 from emt.market_history_updates where "typeID" = $1 and "regionID" = $2 and date_update > $3)`,
		itemID,
		regionID,
		time.Now().AddDate(0,0, daysAvg).Format("2006-01-02"),
	); err != nil {

		return res, err

	}

	if !ex {

		fmt.Println("GetRegionalAverage requesting eve servers", regionID, itemID)

		market, err := esi.MC.MarketHistory(regionID, itemID)
		if err != nil {

			return res, err
		}

		res2 := []MarketHistoryRow{}

		for _, m := range market {

			r := MarketHistoryRow{
				RegionID: regionID,
				TypeID: itemID,
			}

			r.Date = m.Date
			r.OrderCount = m.OrderCount
			r.Volume = m.Volume
			r.Highest = m.Highest
			r.Average = m.Average
			r.Lowest = m.Lowest

			res2 = append(res2, r)

		}

		updDB(res2)
	}

	err = registry.Registry.Connect.SQLX().Get(&res, `select coalesce(avg(average),0) as avg, coalesce(avg(lowest),0) as buy, coalesce(avg(highest),0) as sell from emt.market_history where "typeID" = $1 and "regionID" = $2 and dt > $3`,
		itemID,
		regionID,
		time.Now().AddDate(0, 0, daysAvg).Format("2006-01-02"),
	)

	return

}

func GetMarketHistory(itemID, regionID uint) ([]MarketHistoryRow, error) {

	res := []MarketHistoryRow{}

	// Database has actual data (updated at least 48 hors ago)
	ex := false
	if err := registry.Registry.Connect.SQLX().Get(&ex, `select exists (select 1 from emt.market_history_updates where "typeID" = $1 and "regionID" = $2 and date_update > $3)`,
		itemID,
		regionID,
		time.Now().AddDate(0,0, -2).Format("2006-01-02"),
		); err != nil {

			return res, err

	}

	if ex {

		if err := registry.Registry.Connect.SQLX().Select(&res, `select * from emt.market_history where "typeID" = $1 and "regionID" = $2 and dt > $3`,
			itemID,
			regionID,
			time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
		); err != nil {

			return res, err

		}

		return res, nil

	}

	//Database has outdated entries, must get from esi and update database
	market, err := esi.MC.MarketHistory(regionID, itemID)
	if err != nil {

		return res, err
	}

	//Clean res in case we found something
	res = []MarketHistoryRow{}

	for _, m := range market {

		r := MarketHistoryRow{
			RegionID: regionID,
			TypeID: itemID,
		}

		r.Date = m.Date
		r.OrderCount = m.OrderCount
		r.Volume = m.Volume
		r.Highest = m.Highest
		r.Average = m.Average
		r.Lowest = m.Lowest

		res = append(res, r)

	}

	//Saving to db in a separate goroutine
	go updDB(res)

	return res, nil
}

func updDB (d []MarketHistoryRow) {

	tx, err := registry.Registry.Connect.SQLX().Beginx()
	if err != nil {

		fmt.Println("Cannot start transaction emt.market_history", err.Error())

		return

	}

	var regID, typeID uint

	for _, r := range d {

		regID = r.RegionID
		typeID = r.TypeID

		//Skip today
		if r.Date.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			continue
		}

		ex := false

		if err := tx.Get(&ex, `select exists (select 1 from emt.market_history where "typeID" = $1 and "regionID" = $2 and dt = $3)`,
			r.TypeID,
			r.RegionID,
			r.Date.Format("2006-01-02"),
		); err != nil {

			tx.Rollback()

			fmt.Println("Error selecting from emt.market_history", err.Error(), "\n", r)

			return

		}

		if ex {
			continue
		}

		_, err := tx.Exec(`insert into emt.market_history ("typeID", "regionID", dt, order_count, volume, highest, average, lowest)
					values ($1, $2, $3, $4, $5, $6, $7, $8)`,
			r.TypeID,
			r.RegionID,
			r.Date.Format("2006-01-02"),
			r.OrderCount,
			r.Volume,
			r.Highest,
			r.Average,
			r.Lowest,
		)

		if err != nil {

			tx.Rollback()

			fmt.Println("Error inserting into emt.market_history", err.Error(), "\n", r)

			return

		}

	}

	_, err = tx.Exec(`select 1 from emt.market_history_upsert($1, $2)`, typeID, regID)
	if err != nil {
		tx.Rollback()

		fmt.Println("Error upserting into emt.market_history_updates", err.Error(), "\n", typeID, regID)

		return
	}

	//tx.Exec(`refresh materialized view emt.production_efficiency`)

	tx.Commit()

}