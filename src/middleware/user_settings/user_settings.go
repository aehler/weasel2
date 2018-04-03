package user_settings

import (
	"app"
	"lib/common"
	"lib/auth"
	"lib/personal"
)

func GetUserSettings(c *app.Context) {

	user := c.Get("user").(*auth.User)

	s := common.GetUserSettings(user)

	c.Set("userSettings", s)

	pinned, _ := personal.ListPinned(user)

	c.Set("pinnedBPOCount", len(pinned))

	c.Set("pinnedBPO", pinned)
}