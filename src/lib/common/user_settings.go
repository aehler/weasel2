package common

import (
	"app/registry"
	"lib/auth"
	"encoding/json"
)

func (s *UserSettings) Save(u *auth.User) error {

	return registry.Registry.Session.Upsert(u.SessionID+"_settings", s)

}

func GetUserSettings(u *auth.User) UserSettings  {

	c, err := registry.Registry.Session.Get(u.SessionID+"_settings")
	if err != nil {

		return defaultUserSettings()

	}

	res := UserSettings{}

	err = json.Unmarshal(c, &res)
	if err != nil {

		return defaultUserSettings()

	}

	return res

}

func defaultUserSettings() UserSettings {

	return UserSettings{
		Region: "The Forge",
		RegionID: 10000002,
		MarketShare: 15,
	}

}