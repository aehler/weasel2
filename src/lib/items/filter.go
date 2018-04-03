package items

import (
	"lib/auth"
	"app/registry"
	"fmt"
	"encoding/json"
	"app/form"
)

type Filter struct {
	MaxPrice           float64 `weaselform:"numeric" formLabel:"Price max"`
	MinPrice           float64 `weaselform:"numeric" formLabel:"Price min"`
	MaxCost            float64 `weaselform:"numeric" formLabel:"Cost max"`
	MinCost            float64 `weaselform:"numeric" formLabel:"Cost min"`
	MaxMargin          float64 `weaselform:"numeric" formLabel:"Expected income max"`
	MinMargin          float64 `weaselform:"numeric" formLabel:"Expected income min"`
	MaxRos             float64 `weaselform:"numeric" formLabel:"RoS max"`
	MinRos             float64 `weaselform:"numeric" formLabel:"RoS min"`
	ProductTechLvl     uint    `weaselform:"select" formLabel:"Product Tech level"`
	IsOnMarket         bool    `weaselform:"checkbox" formLabel:"Can be purchased from NPC"`
	ObtainedByResearch bool    `weaselform:"checkbox" formLabel:"Obtained by research"`
	ObtainedByDrop     bool    `weaselform:"checkbox" formLabel:"Drops, LP store etc"`
}

func (p Filter) Options() form.Options {
	return []*form.Option{
		{
			Value : 0,
			Label: "All",
		},
		{
			Value : 1,
			Label: "Tech I",
		},
		{
			Value : 2,
			Label: "Tech II",
		},
		{
			Value : 3,
			Label: "Tech III",
		},
	}
}

func DefaultFilter() *Filter {

	return &Filter{
		IsOnMarket: false,
		ObtainedByResearch: false,
		ObtainedByDrop: false,
		ProductTechLvl: 0,
	}

}

func (f *Filter) UserFilter(u *auth.User) {

	filter, err := registry.Registry.Session.Get(u.SessionID+"_filter")

	if err != nil {

		fmt.Println("UserFilter:", err)

		f = DefaultFilter()

		return

	}

	if err := json.Unmarshal([]byte(filter), f); err != nil {

		fmt.Println(err)

	}

}

func (f *Filter) Save(u *auth.User) {

	if err := registry.Registry.Session.Upsert(u.SessionID+"_filter", f); err != nil {

		fmt.Println("Error saving filters:", err.Error())

	}

}

func (f *Filter) FormatSQL() (sql string) {

	sql = ""

	if f.MaxPrice != 0 {

		sql = fmt.Sprintf("%s and prod_sell_price <= %f", sql, f.MaxPrice)

	}

	if f.MinPrice != 0 {

		sql = fmt.Sprintf("%s and prod_sell_price >= %f", sql, f.MinPrice)

	}

	if f.MaxCost != 0 {

		sql = fmt.Sprintf("%s and material_cost <= %f", sql, f.MaxCost)

	}

	if f.MinCost != 0 {

		sql = fmt.Sprintf("%s and material_cost >= %f", sql, f.MinCost)

	}

	if f.MaxMargin != 0 {

		sql = fmt.Sprintf("%s and net_profit <= %f", sql, f.MaxMargin)

	}

	if f.MinMargin != 0 {

		sql = fmt.Sprintf("%s and net_profit >= %f", sql, f.MinMargin)

	}

	if f.MaxRos != 0 {

		sql = fmt.Sprintf("%s and (ros * 100 - 100) <= %f", sql, f.MaxRos)

	}

	if f.MinRos != 0 {

		sql = fmt.Sprintf("%s and (ros * 100 - 100) >= %f", sql, f.MinRos)

	}

	if f.IsOnMarket {
		sql = fmt.Sprintf("%s and bpo_purchasable = true", sql)
	}

	if f.ObtainedByResearch {
		sql = fmt.Sprintf("%s and obtained_by_research = true", sql)
	}

	if f.ObtainedByDrop {
		sql = fmt.Sprintf("%s and bpo_purchasable = false obtained_by_research = false", sql)
	}

	if f.ProductTechLvl != 0 {
		sql = fmt.Sprintf("%s and tech_lvl = %d", sql, f.ProductTechLvl)
	}

	return
}