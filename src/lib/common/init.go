package common

import (
	"app/registry"
	"strings"
)

type regions struct{
	RegionID uint     `json:"value" db:"regionID"`
	RegionName string `json:"label" db:"regionName"`
}

var regionList []regions

func Init() {

	if err := registry.Registry.Connect.SQLX().Select(&regionList, `select "regionID", "regionName" from evesde."mapRegions" order by "regionName"`); err != nil {
		panic(err.Error())
	}

}

func ListRegions(term string) []regions {

	if term == "" {

		return regionList

	}

	res := []regions{}

	if len(term) <= 2 {
		return res
	}

	for _, reg := range regionList{

		if len(reg.RegionName) >= len(term) && strings.ToLower(reg.RegionName[:len(term)]) == strings.ToLower(term) {

			res = append(res, reg)

		}

	}

	return res

}