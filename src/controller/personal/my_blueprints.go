package personal

import (
	"app"
	"lib/auth"
	"lib/personal"
	"app/crypto"
	"lib/common"
	"lib/items"
	"app/form"
	"app/registry"
	"lib/esi/response"
)

func myBPO(c *app.Context) {

	pinned := c.Get("pinnedBPO").(personal.PinnedBPOs)

	itl := []*items.Item{}

	user := c.Get("user").(*auth.User)

	settings := c.Get("userSettings").(common.UserSettings)

	for _, p := range pinned {

		it, err := items.GetItem(p.TypeID, settings.RegionID, user)

		if err != nil {
			c.RenderHTML("/errors/500.html", map[string]interface {} {
				"Error" : err.Error(),
			})

			c.Stop()

			return
		}

		it.ME = p.ME
		it.TE = p.TE

		itl = append(itl, it)

	}

	genAttrs := personal.ProdLineAttributes{
		LabSlots: 1,
		ManufacturingSlots: 1,
	}

	charSkills := response.CharacterSkills{}

	if err := registry.Registry.Session.Unmarshal(user.SessionID+"_skills", &charSkills); err == nil {

		for _, s := range charSkills.Skills {
			switch s.SkillID {
			case 3387, 24625:
				genAttrs.ManufacturingSlots += s.Active
			case 3406, 24624:
				genAttrs.LabSlots += s.Active
			}

		}

	}

	post := form.New("filters", "to-report", user.SessionID)

	post.Action = "/pinned-report/"

	post.MapStruct(genAttrs)

	post.GetElement("TimeBasis").Options = form.Options{
		{uint(24 * 3), "3 Days"},
		{uint(24 * 7), "Week"},
		{uint(24 * 14), "2 Weeks"},
		{uint(24 * 30), "Month"},
		{uint(24 * 60), "2 Monthes"},
		{uint(24 * 90), "3 Monthes"},
	}

	c.RenderHTML("/my-blueprints.html", map[string]interface {} {
		"items" : itl,
		"form" : post.Context(),
	})

}

func pinnedBPOReport(c *app.Context) {

	user := c.Get("user").(*auth.User)

	settings := c.Get("userSettings").(common.UserSettings)

	post := form.New("filters", "to-report", user.SessionID)

	post.Action = "/pinned-report/"

	genAttrs := personal.ProdLineAttributes{}

	post.MapStruct(genAttrs)

	if c.IsPost() {

		pinned := c.Get("pinnedBPO").(personal.PinnedBPOs)

		post.ParseForm(&genAttrs, c.Request)

		report, err := personal.PinnedReport(user, pinned, settings, genAttrs, &c.Request.Form)

		if err != nil {
			c.RenderHTML("/errors/500.html", map[string]interface {} {
				"Error" : err.Error(),
			})

			c.Stop()

			return
		}

		c.RenderHTML("/my-blueprints-report.html", map[string]interface {} {
			"items" : report.Items,
			"header": report.Header,
			"attrs" : genAttrs,
		})

		return

	}

	c.RenderHTML("/errors/404.html", map[string]interface {} {})

	c.Stop()

	return

}

func togglePinBPO(c *app.Context) {

	user := c.Get("user").(*auth.User)

	c.Request.ParseForm()

	list, err := personal.ListPinned(user)
	if err != nil {

		c.RenderJSON(map[string]interface{}{
			"Error" : err.Error(),
			"Result" : nil,
		})

		c.Stop()

		return

	}

	id, err := crypto.DecryptUrl(c.Request.Form.Get("typeid"))
	if err != nil {

		c.RenderJSON(map[string]interface{}{
			"Error" : err.Error(),
			"Result" : nil,
		})

		c.Stop()

		return
	}

	r, err := list.Toggle(id, user)

	c.RenderJSON(map[string]interface{}{
		"Error" : err,
		"Result" : r,
	})

}