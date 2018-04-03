package index

import (
	"app"
	"app/paginator"
	"lib/items"
	"app/form"
	"lib/auth"
	"lib/common"
)

func Index(c *app.Context) {

	current := paginator.CurrentPage(c)

	user := c.Get("user").(*auth.User)

	settings := c.Get("userSettings").(common.UserSettings)

	limiter := paginator.NewLimiter(20, current)

	sort := "ros"
	sortDir := "desc"
	search := ""
	qs := ""

	if c.GetUrlParam("s") != "" {
		sort = c.GetUrlParam("s")
	}

	if c.GetUrlParam("so") != "" {
		sortDir = c.GetUrlParam("so")
	}

	if c.GetUrlParam("q") != "" {
		search = c.GetUrlParam("q")
		qs = qs + "&q="+search

		filter := items.DefaultFilter()
		filter.Save(user)

	}

	getParams := "/?s="+sort+"&so="+sortDir+"&q="+qs

	post := form.New("filters", "", user.SessionID)

	post.Action = getParams

	filter := items.Filter{}

	filter.UserFilter(user)

	post.MapStruct(filter)

	post.GetElement("ProductTechLvl").Options = filter.Options()

	if c.IsPost() {

		if err := post.ParseForm(&filter, c.Request); err == nil {

			filter.Save(user)

			c.Redirect(getParams)

		} else {

			c.RenderHTML("/errors/500.html", map[string]interface {} {
				"Error" : err.Error(),
			})

			c.Stop()

			return

		}

	}

	itemList, err := items.Index(limiter, sort, sortDir, search, user, settings)

	if err != nil {

		c.RenderHTML("/errors/500.html", map[string]interface {} {
			"Error" : err.Error(),
		})

		c.Stop()

		return

	}

	c.RenderHTML("/index.html", map[string]interface {} {
		"items"   : itemList.Items,
		"sort"    : sort,
		"sortDir" : sortDir,
		"search"  : search,
		"qs"      : qs,
		"totals"  : itemList.Totals,
		"form"    : post.Context(),
		"paginator": paginator.NewPaginator(
			current,
			itemList.Totals.Total,
			limiter.Limit(),
			"/",
			map[string]string{
				"q" : search,
				"s" : sort,
				"so" : sortDir,
			},
		),
	})

}