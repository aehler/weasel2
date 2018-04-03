package items

import (
	"app/registry"
	"time"
	"encoding/json"
	"errors"
	"math"
	"fmt"
	"strings"
	"lib/auth"
	"lib/esi/response"
)

type Type struct {
	TypeID uint `db:"typeID"`
	TypeName string `db:"typeName"`
	BpoName string `db:"-"`
}

type Item struct{
	TypeID uint `db:"bpo_type_id"`
	ProductTypeID uint `db:"product_type_id"`
	BpoName string `db:"bpo_name"`
	ProductName string `db:"product_name"`
	MetaTypeName string `db:"meta_type_name"`
	TechLvl *uint `db:"tech_lvl"`
	Volume float64 `db:"volume"`
	Batch uint `db:"batch"`
	SalesAvg float64 `db:"sales_avg"`
	ManufTime time.Duration `db:"manufacturing_time"`
	RequiredSkills Skills `db:"skills"`
	BPOPurchasable bool `db:"bpo_purchasable"`
	MarketShare uint `db:"market_share"`
	Materials Mats
	TimeModifiers float64
	AdjTime time.Duration
	TotalCostBuy float64
	TotalCostProduce float64
	TotalVolumeBuy float64
	TotalVolumeProduce float64
	MarginalIncomeBuy float64
	MarginalIncomeProduce float64
	IskhBuy float64
	IskhProduce float64
	ResearchData *Research
	ME uint
	TE uint
}

type Research struct {
	Quantity uint `db:"quantity"`
	Probability float64 `db:"probability"`
	AdjProbability float64 `db:"-"`
	Materials ResearchMats `db:"datacores"`
	ResearchTime time.Duration `db:"research_time"`
	TimeModifiers float64 ``
	AdjTime time.Duration
	RequiredSkills Skills `db:"skills"`
}

type ResearchMats []ResearchMat

type ResearchMat struct {
	TypeID uint `json:"mat_type_id"`
	Quantity uint `json:"quantity"`
	MaterialTypeName string `json:"mat_name"`
	Cost float64 `json:"cost"`
}


func GetType(id uint) (res Type, err error) {

	err = registry.Registry.Connect.SQLX().Get(&res, `select "typeID", "typeName" from evesde."invTypes" where "typeID" = $1`, id)

	return

}

func GetTypeByBPO(id uint) (res Type, err error) {

	err = registry.Registry.Connect.SQLX().Get(&res, `select tpp."typeID", tpp."typeName" from evesde."invTypes" as bp
	left join evesde."industryActivityProducts" as ap using("typeID")
	left join evesde."invTypes" as tpp on tpp."typeID" = ap."productTypeID"
	where bp."typeID" = $1 and ap."activityID" = $2`, id, ActivityManufacturing)

	return

}

func GetItem(id, regionId uint, user *auth.User, bpoattr ...uint) (*Item, error) {

	cacheKey := fmt.Sprintf("i%dr%d", id, regionId)

	res := GetFromCache(cacheKey)

	if res != nil {

		fmt.Println("From cache")

		return res, nil
	} else {
		res = &Item{}
	}

	// fmt.Println(id, strings.Join(timeBonusAttributes, ","))

	if err := registry.Registry.Connect.SQLX().Get(res, `select tp."typeID" as bpo_type_id, tpp."typeID" as product_type_id, tp."typeName" as bpo_name,
tpp."typeName" as product_name,
case when pe."metaGroupID" is not null then pe."metaGroupName"
	when pe."metaGroupID" is null and pe.tech_lvl is not null then
		case pe.tech_lvl when 1 then 'Tech I'
		when 2 then 'Tech II'
		when 3 then 'Tech III'
		else 'unknown' end
	else 'unknown'
end as meta_type_name,
pe.bpo_purchasable,
prod.quantity as batch, tpp.volume, ma.avg_price as sales_avg,
(a.time::bigint * 1000000000) as manufacturing_time,
	coalesce((
	select json_agg(json_build_object('skillID', skills."skillID", 'lvl', skills.level, 'name', ts."typeName", 'treelvl', 0,
		'time_bonus', (select sum(coalesce(skill2."valueInt", skill2."valueFloat")) FROM evesde."dgmTypeAttributes" AS skill2
			where skill2."typeID" = skills."skillID" and skill2."attributeID" = any($3))
	))
		from evesde."industryActivitySkills" as skills
		left join evesde."invTypes" as ts on skills."skillID" = ts."typeID"
		where skills."typeID" = tp."typeID" and skills."activityID" = $2
	), '[]') as skills,
	pe.tech_lvl,
	coalesce(ag.volume, 0) as market_share
from evesde."industryBlueprints" as bpo
left join evesde."invTypes" as tp using("typeID")
left join evesde."invGroups" as g using("groupID")
left join evesde."industryActivityProducts" as prod using("typeID")
left join evesde."invTypes" as tpp on prod."productTypeID" = tpp."typeID"
left join emt.market_avg as ma on ma.type_id = tpp."typeID"
left join emt.production_efficiency as pe on pe.product_type_id = tpp."typeID"
left join evesde."invMarketGroups" as mg on mg."marketGroupID" = tpp."marketGroupID"
left join evesde."industryActivity" as a on a."typeID" = tp."typeID" and a."activityID" = $2
left join emt.avg_grouped as ag on ag."typeID" = tpp."typeID" and ag."regionID" = $4
where tp."typeID" = $1
and prod."activityID" = $2`, id, ActivityManufacturing, fmt.Sprintf("{%s}", strings.Join(timeBonusAttributes, ",")), regionId); err != nil {

		fmt.Println("Main query error:", err.Error())

		return res, err

	}

	teBonuses := map[uint]int{}
	teBonusesLevels := map[uint]uint{}
	charSkills := response.CharacterSkills{}
	hs := map[uint]uint{}

	if res.TechLvl != nil {
		res.ME = DefaultBPOAttr[*res.TechLvl][0]
		res.TE = DefaultBPOAttr[*res.TechLvl][0]

		if !res.BPOPurchasable && *res.TechLvl == 1 {
			res.ME = 0
			res.TE = 0
		}
	}

	switch len(bpoattr) {
	case 1:
		res.ME = bpoattr[0]

	case 2:
		res.ME = bpoattr[0]
		res.TE = bpoattr[1]

	}

	registry.Registry.Session.Unmarshal(user.SessionID+"_skills", &charSkills)

	for _, ownSkill := range charSkills.Skills{

		hs[ownSkill.SkillID] = ownSkill.Active

	}

	for i, _ := range res.RequiredSkills {

		if err := (res.RequiredSkills[i]).SubSkills(timeBonusAttributes); err != nil {

			return res, err

		}

	}

	for _, s := range res.RequiredSkills {

		if s.TimeBonus < 0 {
			teBonuses[s.SkillID] = s.TimeBonus
		}

		for _, ss := range s.Tree {

			if ss.TimeBonus < 0 {
				teBonuses[ss.ReqID] = ss.TimeBonus
			}

		}
	}

	for _, s := range res.RequiredSkills {

		for sid, _ := range teBonuses {

			if s.SkillID == sid {

				ownSkillActive, ok := hs[sid]
				if !ok {

					ownSkillActive = s.Level

				} else {

					if s.Level > ownSkillActive {

						ownSkillActive = s.Level

					}
				}

				if v, ok := teBonusesLevels[s.SkillID]; ok {

					if v < s.Level {

						teBonusesLevels[s.SkillID] = ownSkillActive

					}

				} else {

					teBonusesLevels[s.SkillID] = ownSkillActive

				}

			}

		}

		for _, ss := range s.Tree {

			for sid, _ := range teBonuses {

				if ss.ReqID == sid {

					ownSkillActive, ok := hs[sid]
					if !ok {

						ownSkillActive = ss.Level

					} else {

						if s.Level > ownSkillActive {

							ownSkillActive = ss.Level

						}
					}

					if v, ok := teBonusesLevels[ss.ReqID]; ok {

						if v < ss.Level {

							teBonusesLevels[ss.ReqID] = ownSkillActive

						}

					} else {

						teBonusesLevels[ss.ReqID] = ownSkillActive

					}

				}

			}

		}
	}

	for k, s := range hs {

		if k == IndustryID || k == AdvIndustryID {

			if _, ok := teBonusesLevels[k]; !ok {
				teBonusesLevels[k] = s
			}

		}

	}

	res.SalesAvg = res.SalesAvg * float64(res.Batch)

	res.TimeModifiers = 1 - float64(res.TE / 100)

	for sid, lvl := range teBonusesLevels {

		res.TimeModifiers = res.TimeModifiers * (1 + float64(lvl) * (float64(teBonuses[sid])/100))

	}

	res.AdjTime = time.Duration(math.Floor(res.TimeModifiers * float64(res.ManufTime)/1000000000)) * time.Second

	if err := res.GetMaterials(); err != nil {

		return res, err

	}

	for _, mat := range res.Materials {

		if mat.Lvl == 1 {

			res.TotalCostBuy = res.TotalCostBuy + mat.Cost
			res.TotalVolumeBuy = res.TotalVolumeBuy + mat.Volume

			if mat.ComponentBPO != 0 {
				res.TotalVolumeProduce = res.TotalVolumeProduce + mat.SubmatsTotalVolume
				res.TotalCostProduce = res.TotalCostProduce + mat.SubmatsTotalCost
			} else {
				res.TotalVolumeProduce = res.TotalVolumeProduce + mat.Volume
				res.TotalCostProduce = res.TotalCostProduce + mat.Cost
			}
		}

	}

	res.MarginalIncomeBuy = res.SalesAvg - res.TotalCostBuy
	res.MarginalIncomeProduce = res.SalesAvg - res.TotalCostProduce
	res.IskhBuy = res.MarginalIncomeBuy / (res.ManufTime.Seconds() / 60 / 60)

	if res.TechLvl != nil && (*res.TechLvl == 2 || *res.TechLvl == 3) {

		res.GetResearch()

	}

	res.Cache(cacheKey)

	return res, nil
}

func (m *Item) Cache(key string) {

	if b, err := json.Marshal(m); err == nil {

		registry.Registry.Session.Cache(key, b, time.Second * 15)

	} else {

		fmt.Println("failed to store item in cache, key:", key)

	}

}

func GetFromCache(key string) *Item {

	if b, err := registry.Registry.Session.GetNoTouch(key); err == nil {

		res := Item{}

		if err := json.Unmarshal(b, &res); err == nil {

			return &res

		}

	}

	return nil
}

func (m *Item) GetResearch() error {

	res := Research{}

	if err := registry.Registry.Connect.SQLX().Get(&res, `select act.quantity, prob.probability, json_agg(json_build_object('mat_type_id', mats."materialTypeID", 'quantity', mats.quantity, 'mat_name', dctp."typeName", 'cost', (ma.avg_price * act.quantity))) as datacores,
(ia.time::bigint * 1000000000) as research_time,
coalesce((
	select json_agg(json_build_object('skillID', skills."skillID", 'lvl', skills.level, 'name', ts."typeName", 'treelvl', 0,
		'time_bonus', (select sum(coalesce(skill2."valueInt", skill2."valueFloat")) FROM evesde."dgmTypeAttributes" AS skill2
			where skill2."typeID" = skills."skillID" and skill2."attributeID" = any($3)),
		'prob_bonus', (select sum(coalesce(skill2."valueInt", skill2."valueFloat")) FROM evesde."dgmTypeAttributes" AS skill2
			where skill2."typeID" = skills."skillID" and skill2."attributeID" = $4)
	))
		from evesde."industryActivitySkills" as skills
		left join evesde."invTypes" as ts on skills."skillID" = ts."typeID"
		where skills."typeID" = ia."typeID" and skills."activityID" = $2
	), '[]') as skills
from evesde."industryActivityProducts" as act
left join evesde."industryActivity" as ia using ("typeID")
left join evesde."industryActivityProbabilities" as prob using ("typeID")
left join evesde."industryActivityMaterials" as mats using ("typeID")
left join evesde."invTypes" as dctp on dctp."typeID" = mats."materialTypeID"
left join emt.market_avg as ma on ma.type_id = dctp."typeID"
where act."productTypeID" = $1
and act."activityID" = $2
and mats."activityID" = $2
and ia."activityID" = $2
group by act.quantity, prob.probability, ia.time, ia."typeID"
		`,
		m.TypeID,
		ActivityInvention,
		fmt.Sprintf("{%s}", strings.Join(invTimeBonusAttributes, ",")),
		inventionPropabilityMultiplier,
	); err != nil {

		fmt.Println("Error getting research data", err.Error())

		return err
	}

	teBonuses := map[uint]int{}
	teBonusesLevels := map[uint]uint{}
	var rem uint
	var sc uint

	for i, _ := range res.RequiredSkills {

		if err := (&res.RequiredSkills[i]).SubSkills(invTimeBonusAttributes); err != nil {

			return err

		}

	}

	for _, s := range res.RequiredSkills {

		isrem := false

		for _, enc := range encMethodsIDS {
			if s.SkillID == enc {

				rem = s.Level

				isrem = true

			}
		}

		if !isrem {
			sc = sc + s.Level
		}

		if s.TimeBonus < 0 {
			teBonuses[s.SkillID] = s.TimeBonus
		}

		for _, ss := range s.Tree {

			if ss.TimeBonus < 0 {
				teBonuses[ss.ReqID] = ss.TimeBonus
			}

		}
	}

	for _, s := range res.RequiredSkills {

		for sid, _ := range teBonuses {

			if s.SkillID == sid {

				if v, ok := teBonusesLevels[s.SkillID]; ok {

					if v < s.Level {

						teBonusesLevels[s.SkillID] = s.Level

					}

				} else {

					teBonusesLevels[s.SkillID] = s.Level

				}

			}

		}

		for _, ss := range s.Tree {

			for sid, _ := range teBonuses {

				if ss.ReqID == sid {

					if v, ok := teBonusesLevels[ss.ReqID]; ok {

						if v < ss.Level {

							teBonusesLevels[ss.ReqID] = ss.Level

						}

					} else {

						teBonusesLevels[ss.ReqID] = ss.Level

					}

				}

			}

		}
	}

	res.TimeModifiers = 1

	res.AdjProbability = res.Probability * ( 1 + (float64(rem)/40) + (float64(sc)/30) )

	for sid, lvl := range teBonusesLevels {

		res.TimeModifiers = res.TimeModifiers * (1 + float64(lvl) * (float64(teBonuses[sid])/100))

	}

	res.AdjTime = time.Duration(math.Floor(res.TimeModifiers * float64(res.ResearchTime)/1000000000)) * time.Second

	m.ResearchData = &res

	return nil
}

func (u *ResearchMats) Scan(src interface{}) error {

	var source []byte

	switch src.(type) {

	case string:

		source = []byte(src.(string))

	case []byte:

		source = src.([]byte)

	default:

		return errors.New("Incompatible type for ResearchMats")
	}

	if err := json.Unmarshal(source, &u); err != nil {

		return err
	}

	return nil
}