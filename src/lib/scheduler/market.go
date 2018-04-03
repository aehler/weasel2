package scheduler

import (
	"fmt"
	"lib/esi"
	"app/registry"
	"lib/market"
	"time"
)

func updMarketData() {

	fmt.Println("Updating market prices")

	mp, err := esi.MC.MarketPrices()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	tx, _ := registry.Registry.Connect.SQLX().Beginx()

	if _, err := tx.Exec(`delete from emt.market_avg`); err != nil {

		fmt.Println(err.Error())

		tx.Rollback()

		return

	}

	for _, m := range mp {

		if _, err := tx.Exec(`insert into emt.market_avg (type_id, avg_price, adj_price) values ($1, $2, $3)`,
			m.TypeID,
			m.AvgPrice,
			m.AdjPrice,
		); err != nil {

			fmt.Println(err.Error())

			tx.Rollback()

			return
		}

	}

	tx.Commit()

	fmt.Println("Done")

	return
}

func updMarketAverages() {

	fmt.Println("Updating market averages")

	if _, err := registry.Registry.Connect.SQLX().Exec(`refresh materialized view emt.avg_grouped`); err != nil {

		fmt.Println(err.Error())

		return

	}

}

func updRegionalMarket() {

	fmt.Println("Updating regional data")

	itemIds := []uint{}

	if err := registry.Registry.Connect.SQLX().Select(&itemIds, `select product_type_id from emt.production_efficiency`); err != nil {

		fmt.Println(err.Error())

		return

	}

	go func(){

		for _, i := range itemIds {

			_, err := market.GetMarketHistory(i, 10000043)
			if err != nil {

				fmt.Println("Error getting history", err.Error())

			}

			time.Sleep(5 * time.Second)

		}

	}()

}