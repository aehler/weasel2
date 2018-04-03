package paginator

import (
	"app"
	"strconv"
)

func CurrentPage(c *app.Context) uint {

	if page, err := strconv.ParseUint(c.Param("page"), 10, 0); err == nil {

		return uint(page)
	}

	return 1
}
