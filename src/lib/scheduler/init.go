package scheduler

import (
	"github.com/revel/cron"
	"fmt"
)

func init() {

	c := cron.New()

	fmt.Println("Added", len(c.Entries()), "cron jobs")

	c.Start()
}

