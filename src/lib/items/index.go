package items

import (
	"app/registry"
	"app/paginator"
	"fmt"
	"errors"
	"strings"
	"lib/auth"
	"encoding/json"
	"lib/common"
)

type ItemIndex struct {
	TypeName string `db:"product"`
	GroupName string `db:"prod_group_name"`
	TypeID uint `db:"typeID"`
	ProductTypeID uint `db:"product_type_id"`
	AvgPrice float64 `db:"prod_sell_price"`
	AvgCost float64 `db:"material_cost"`
	Income float64 `db:"net_profit"`
	RoS float64 `db:"ros"`
	MetaGroupName *string `db:"metaGroupName"`
	MetaGroupID *uint `db:"metaGroupID"`
	TechLvl *uint `db:"tech_lvl"`
	MetaLvl *uint `db:"meta_lvl"`
	BPOPurchasable bool `db:"bpo_purchasable"`
	NotPrecise bool `db:"not_precise"`
	Mats SimpleMats `db:"mats_json"`
	OrdersVol uint `db:"orders_volume"`
}

type SimpleMats []SimpleMat

type SimpleMat struct {
	TypeID uint `json:"material_type_id"`
	Material string `json:"material"`
	Price float64 `json:"avg_price"`
	Quantity float64 `json:"quantity"`
}

func (u *SimpleMats) Scan(src interface{}) error {

	var source []byte

	switch src.(type) {

	case string:

		source = []byte(src.(string))

	case []byte:

		source = src.([]byte)

	default:

		return errors.New("Incompatible type for SimpleMats")
	}

	if err := json.Unmarshal(source, &u); err != nil {

		return err
	}

	return nil
}

type ItemsIndex struct {
	Items []ItemIndex
	Totals GroupValues
}

type GroupValues struct {
	Total uint
	MaxPrice float64
	MaxRoS float64
	MaxProfit float64
}

func Index(limiter paginator.Limiter, sort, sortDir, search string, user *auth.User, settings common.UserSettings) (ItemsIndex, error) {

	res := ItemsIndex{}

	filter := &Filter{}

	filter.UserFilter(user)

	if err := registry.Registry.Connect.SQLX().Get(&res.Totals, fmt.Sprintf(`select count("typeID") as total, coalesce(max(prod_sell_price),0) as maxprice,
	coalesce(max(ros * 100 - 100),0) as maxros, coalesce(max(prod_sell_price - material_cost),0) as maxprofit
	from emt.production_efficiency where material_cost is not null
	and prod_sell_price != 0
	%s
	and case when $1 != '' then to_tsquery($1) @@ norm_vc else true end`, filter.FormatSQL()), strings.Replace(strings.TrimSpace(search), " ", "+", -1)); err != nil {

		return res, err

	}

	switch sort{

	case "ros":
	case "sa":
		sort = "prod_sell_price"
	case "ca":
		sort = "material_cost"
	case "inc":
		sort = "net_profit"
	case "ov":
		sort = "orders_volume"

	default:
		return res, errors.New("Invalid params")

	}

	switch sortDir{

	case "asc", "desc":

	default:
		return res, errors.New("Invalid params")

	}

	if err := registry.Registry.Connect.SQLX().Select(&res.Items, fmt.Sprintf(`select net_profit, ros * 100 - 100 as ros, product, pe."typeID", product_type_id,
		prod_sell_price, material_cost, prod_group_name,
		"metaGroupID", "metaGroupName", bpo_purchasable, tech_lvl, meta_lvl, (select bool_or(a is null) from unnest(avg_prices) s(a)) as not_precise, mats_json,
		case coalesce(t.days_back, 0) when 0 then
		0
		else
		coalesce(round((t.volume * t.days)/t.days_back), 0)
		end as orders_volume
		from emt.production_efficiency as pe
		left join emt.avg_grouped as t on t."typeID" = pe.product_type_id and t."regionID" = $4
		where material_cost is not null and prod_sell_price != 0 and prod_sell_price is not null
		%s
		and case when $3 != '' then to_tsquery($3) @@ norm_vc else true end
		order by %s %s
		limit $1 offset $2`, filter.FormatSQL(), sort, sortDir), limiter.Limit(), limiter.Offset(), strings.Replace(strings.TrimSpace(search), " ", "+", -1),
			settings.RegionID,
			); err != nil {

		return res, err

	}

	return res, nil

}