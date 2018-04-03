package scheduler

import (
	"fmt"
	"app/registry"
	"lib/items"
)

func updAvgByItems() {

	fmt.Println("Updating local data by average prices")

	_, err := registry.Registry.Connect.SQLX().Exec(`refresh materialized view emt.production_efficiency`)

	if err != nil {

		fmt.Println(err)

	}

	err = items.GetDecryptorData()

	if err != nil {

		fmt.Println(err)

	}

	fmt.Println("Done")
}