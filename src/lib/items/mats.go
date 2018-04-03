package items

import (
	"app/registry"
	"lib/market"
	"fmt"
	"strings"
)

type Mats []Mat

type Mat struct {
	BpoID uint `db:"bpo_id"`
	MaterialID uint `db:"material_id"`
	Quantity float64 `db:"quantity"`
	Batch uint `db:"batch"`
	MaterialName string `db:"typeName"`
	AvgPrice float64 `db:"avg_price"`
	Cost float64 `db:"material_cost"`
	Volume float64 `db:"volume"`
	TotalVolume float64 `db:"total_vol"`
	ComponentBPO uint `db:"component_bpo_id"`
	Lvl uint `db:"lvl"`
	ActivityID uint `db:"mats_activity"`
	ActivityName string `db:"activityName"`
	LPath Lpath `db:"l_path"`
	SubmatsTotalVolume float64 `db:"-"`
	SubmatsTotalCost float64 `db:"-"`
	SellPrice float64
	BuyPrice float64
	ProdPrice float64
	BetterBSP string
}

type Lpath string

func (l Lpath) In(parentPath Lpath) bool {

	return len(string(parentPath)) <= len(string(l)) &&
		strings.Compare(string(parentPath), string(l[0:len(parentPath)])) == 0

}

func (m *Item) GetMaterials() error {

	res := Mats{}

	fmt.Println(m.TypeID)

	if err := registry.Registry.Connect.SQLX().Select(&res, `with recursive r as (

	select bpo."typeID" as bpo_id, mat."typeID" as material_id, mats.quantity, mat."typeName", coalesce(ma.avg_price, 0) as avg_price, coalesce(ma.avg_price * mats.quantity, 0) as material_cost, mat.volume, mat.volume * mats.quantity as total_vol
	,coalesce(mats_bpo."typeID", 0) as component_bpo_id, 1 as lvl, mats."activityID" as mats_activity, act_c."activityName",
	(select prod.quantity from evesde."industryActivityProducts" as prod where prod."typeID" = bpo."typeID" and prod."activityID" = 1) as batch,
	coalesce(mats_bpo."typeID", 0)::text as l_path
	from evesde."industryBlueprints" as bpo
	left join evesde."industryActivityMaterials" as mats on mats."typeID" = bpo."typeID"
	left join evesde."invTypes" as mat on mat."typeID" = mats."materialTypeID"
	left join emt.market_avg as ma on ma.type_id = mats."materialTypeID"
	left join evesde."ramActivities" as act_c using ("activityID")
	left join evesde."industryActivityProducts" as mats_bpo on mats_bpo."productTypeID" = mat."typeID" and mats_bpo."activityID" = 1
	where bpo."typeID" = $1
	and mats."activityID" = $2

	union

	select bpo."typeID" as bpo_id, mat."typeID" as material_id, mats.quantity, mat."typeName", coalesce(ma.avg_price, 0) as avg_price, coalesce(ma.avg_price * mats.quantity, 0) as material_cost, mat.volume, mat.volume * mats.quantity as total_vol
	,coalesce(mats_bpo."typeID", 0) as component_bpo_id, lvl+1 as lvl, mats."activityID" as mats_activity, act_c."activityName",
	(select prod.quantity from evesde."industryActivityProducts" as prod where prod."typeID" = bpo."typeID" and prod."activityID" = 1) as batch,
	l_path || '.' || coalesce(mats_bpo."typeID", 0) as l_path
	from evesde."industryBlueprints" as bpo
	left join evesde."industryActivityMaterials" as mats on mats."typeID" = bpo."typeID"
	left join evesde."invTypes" as mat on mat."typeID" = mats."materialTypeID"
	left join emt.market_avg as ma on ma.type_id = mats."materialTypeID"
	left join evesde."ramActivities" as act_c using ("activityID")
	left join evesde."industryActivityProducts" as mats_bpo on mats_bpo."productTypeID" = mat."typeID" and mats_bpo."activityID" = 1
	join r on bpo."typeID" = r.component_bpo_id
	and mats."activityID" = $2
)
select * from r`, m.TypeID, ActivityManufacturing); err != nil {

		return err

	}

	for i, mat := range res {

		if mat.ComponentBPO != 0 && mat.Lvl == 1 {

			res[i].SubmatsTotalVolume, res[i].SubmatsTotalCost = res.prepareRecursive(mat.ComponentBPO, mat.LPath)

		}

	}

	m.Materials = res

	return nil

}

func (m *Mats) PrepareRecursiveRegional(componentID, regionID uint, componentLPath Lpath, ttt int, plannedBatch float64) (smv, smc float64) {

	var currentComponent *Mat

	for _, mat := range *m {

		if mat.ComponentBPO == componentID && mat.LPath.In(componentLPath) {

			currentComponent = &mat

			break
		}

	}

	if currentComponent == nil {
		return
	}

	for i, mat := range *m {

		if mat.BpoID == currentComponent.ComponentBPO && mat.LPath.In(componentLPath) {

			rm, err := market.GetRegionalAverage(mat.MaterialID, regionID, ttt * -1)
			if err != nil {
				return smv, smc
			}

			(*m)[i].SellPrice = rm.Sell * mat.Quantity / float64(mat.Batch)
			(*m)[i].BuyPrice = rm.Buy * mat.Quantity / float64(mat.Batch)

			smv = smv + (*m)[i].SellPrice
			smc = smc + (*m)[i].BuyPrice

			if mat.ComponentBPO != 0 {

				(*m)[i].ProdPrice, _ = m.PrepareRecursiveRegional(mat.ComponentBPO, regionID, mat.LPath, ttt, plannedBatch)

			}

		}

	}

	return smv, smc

}

func (m *Mats) prepareRecursive(componentID uint, componentLPath Lpath) (smv, smc float64) {

	var currentComponent *Mat

	for _, mat := range *m {

		if mat.ComponentBPO == componentID && mat.LPath.In(componentLPath) {

			currentComponent = &mat

			break
		}

	}

	if currentComponent == nil {
		return
	}

	for i, mat := range *m {

		//fmt.Println(currentComponent.LPath, mat.LPath, len(currentComponent.LPath))

		//if len(currentComponent.LPath) <= len(mat.LPath) && strings.Compare(currentComponent.LPath, mat.LPath[0:len(currentComponent.LPath)]) == 0 {
		//	if strings.Index(mat.MaterialName, "Photon") != -1 {
		//		fmt.Println(mat.MaterialName, currentComponent.MaterialName, mat.Quantity, currentComponent.Quantity, mat.LPath, currentComponent.LPath)
		//	}
		//}

		if mat.BpoID == currentComponent.ComponentBPO && mat.LPath.In(currentComponent.LPath) {

			(*m)[i].Quantity = mat.Quantity / float64(mat.Batch) * currentComponent.Quantity
			(*m)[i].Cost = mat.Cost / float64(mat.Batch) * float64(currentComponent.Quantity)
			(*m)[i].TotalVolume = mat.TotalVolume / float64(mat.Batch) * float64(currentComponent.Quantity)

			smv = smv + (*m)[i].TotalVolume
			smc = smc + (*m)[i].Cost

			if mat.ComponentBPO != 0 {

				(*m)[i].SubmatsTotalVolume, (*m)[i].SubmatsTotalCost = m.prepareRecursive(mat.ComponentBPO, mat.LPath)

			}

		}

	}

	return smv, smc

}