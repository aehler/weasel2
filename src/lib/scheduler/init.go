package scheduler

import (
	"github.com/revel/cron"
	"fmt"
)

func init() {

	c := cron.New()

	c.AddFunc("0 45 * * * *", updMarketData)

	c.AddFunc("0 50 * * * *", updAvgByItems)

	c.AddFunc("0 */5 * * * *", updMarketAverages)

	c.AddFunc("0 48 17 * * 4", updRegionalMarket) //At 17:48 every Thursday

	fmt.Println("Added", len(c.Entries()), "cron jobs")

	c.Start()
}

