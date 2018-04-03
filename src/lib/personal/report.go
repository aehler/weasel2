package personal

import (
	"lib/auth"
	"lib/items"
	"lib/common"
	"strconv"
	"lib/market"
	"net/url"
	"app/crypto"
	"helper"
	"math"
	"fmt"
)

type PReport struct {
	Header h
	Items []ReportItem
}

type ReportItem struct {
	Item *items.Item
	SalesBatch uint
	RegionalAvg float64
	RawProfit float64
	MaterialCost float64
	BasedOnShare bool
	PlannedBatch uint
	TimeToSell uint
	InventionIncluded bool
	InventionCostAdj float64
	InventionCostAdjWithDecryptor float64
	InventionCostAdjWithDecryptorTotal float64
	InventionChanceWithDecryptor float64
	InventionJobRuns uint
	MaxRuns uint
	JobRuns uint
	InventionCosts float64
	Decryptor items.Decriptor
	ME uint
	TE uint
}

type h struct {
	Basis string
}

const regAvgPeriod = 30

type decr struct{
	d items.Decriptor
	invCost float64
	mCost float64
	q float64
}

func selectOptimalDecryptor(invCost float64, it items.Item, settings common.UserSettings) (items.Decriptor, error) {

	res := []decr{}

	var mc float64

	decryptors := append([]items.Decriptor{items.Decriptor{
		InvChance: 1,
		MaxRuns: 0,
		AdjPrice: 0,
		TypeID: 0,
		TypeName: "No decryptor",
	}}, items.Decryptors()...)

	for _, d := range decryptors {

		me := 2 + d.ME
		var matcost, q float64

		for _, m := range it.Materials {

			if m.Lvl == 1 {

				m.Quantity = helper.Round(float64(m.Quantity * float64(d.MaxRuns + int(it.ResearchData.Quantity)) * float64(1 - float64(me) / 100)))
				q = m.Quantity

				rm, err := market.GetRegionalAverage(m.MaterialID, settings.RegionID, regAvgPeriod * -1)
				if err != nil {
					return items.Decriptor{}, err
				}

				matcost = m.Quantity * rm.Avg

			}

		}

		res = append(res, decr{
			d : d,
			invCost: (invCost + d.AdjPrice) / (float64(int(it.ResearchData.Quantity) + d.MaxRuns) * (it.ResearchData.AdjProbability * d.InvChance) ),
			mCost: matcost,
			q : q,
		})

	}

	mc = res[0].mCost + res[0].invCost
	cd := res[0].d

	for _, r := range res {

		if mc > (r.mCost + r.invCost) / r.q {

			mc = (r.mCost + r.invCost) / r.q

			cd = r.d

		}

	}

	return cd, nil

}

func PinnedReport(user *auth.User, pinned PinnedBPOs, settings common.UserSettings, genAttrs ProdLineAttributes, form *url.Values) (PReport, error) {

	res := PReport{}

	ttt, err := strconv.Atoi(genAttrs.TimeBasis)
	if err != nil {
		ttt = 168
	}

	ttt = ttt / 24

	switch ttt {
	case 3:
		res.Header.Basis = "three days"
	case 7:
		res.Header.Basis = "a week"
	case 14:
		res.Header.Basis = "two weeks"
	case 28:
		res.Header.Basis = "a month"
	case 28 * 2:
		res.Header.Basis = "two monthes"
	case 28 * 3:
		res.Header.Basis = "three monthes"
	}

	for _, p := range pinned {

		//Production plan part
		eid := crypto.EncryptUrl(p.TypeID)

		it, err := items.GetItem(p.TypeID, settings.RegionID, user)
		if err != nil {
			return res, err
		}

		plannedBatch, err := strconv.Atoi(form.Get("batch_"+eid))
		if err != nil || plannedBatch <= 0 {
			plannedBatch = int(it.MarketShare * uint(ttt) * settings.MarketShare / 100)
		}

		ra, err := market.GetRegionalAverage(it.ProductTypeID, settings.RegionID, -90)
		if err != nil {
			return res, err
		}

		r := ReportItem{
			Item: it,
			SalesBatch: uint(it.MarketShare * uint(ttt) * settings.MarketShare / 100),
			RegionalAvg: ra.Sell,
			PlannedBatch : uint(plannedBatch),
		}

		r.BasedOnShare = int(r.SalesBatch) == plannedBatch

		if (it.MarketShare * settings.MarketShare / 100) != 0 {
			r.TimeToSell = uint(plannedBatch) / (it.MarketShare * settings.MarketShare / 100)
		} else if it.MarketShare == 0 {
			r.TimeToSell = 999
		} else {
			r.TimeToSell = 999
		}

		invention := form.Get("invention_"+eid)
		if invention == "on" {
			r.InventionIncluded = true
		}

		bme, err := strconv.Atoi(form.Get("me_"+eid))
		if err != nil {
			r.ME = it.ME
		}
		r.ME = uint(bme)

		r.MaxRuns = 999999
		r.JobRuns = 1

		//Invention part
		if r.InventionIncluded {

			var invCost float64 = 0

			for i, research := range it.ResearchData.Materials {

				rm, err := market.GetRegionalAverage(research.TypeID, settings.RegionID, regAvgPeriod * -1)
				if err != nil {
					return res, err
				}

				it.ResearchData.Materials[i].Cost = rm.Sell * float64(research.Quantity)

				invCost = invCost + rm.Sell * float64(research.Quantity)

			}

			costPerRunAdjChance := invCost / (float64(it.ResearchData.Quantity) * it.ResearchData.AdjProbability)

			r.InventionCostAdj = costPerRunAdjChance

			recommendedDecryptor, err := selectOptimalDecryptor(invCost, *it, settings)
			if err != nil {
				return res, err
			}

			r.Decryptor = recommendedDecryptor

			r.ME = uint(2 + r.Decryptor.ME)

			wdpmin := (invCost + r.Decryptor.AdjPrice) / (float64(int(it.ResearchData.Quantity) + r.Decryptor.MaxRuns) * (it.ResearchData.AdjProbability * r.Decryptor.InvChance) )

			r.InventionCostAdjWithDecryptor = wdpmin

			r.InventionCostAdjWithDecryptorTotal = float64(r.SalesBatch) * wdpmin

			r.InventionJobRuns = uint(float64(r.SalesBatch / it.ResearchData.Quantity) / it.ResearchData.Probability)

			r.MaxRuns = uint(r.Decryptor.MaxRuns + int(it.ResearchData.Quantity))

			if r.Decryptor.TypeID != 0 {

				r.InventionChanceWithDecryptor = it.ResearchData.AdjProbability * r.Decryptor.InvChance

				//									Plan to produce                    max runs with decryptor                       probability with decryptor
				r.InventionJobRuns = uint(math.Floor(float64(r.SalesBatch / uint((int(it.ResearchData.Quantity) + r.Decryptor.MaxRuns))) / r.InventionChanceWithDecryptor))

				r.InventionCosts = r.Decryptor.AdjPrice * float64(r.InventionJobRuns)

			}

			for _, research := range it.ResearchData.Materials {

				r.InventionCosts = r.InventionCosts + research.Cost * float64(r.InventionJobRuns)

			}

		}

		copyRuns := int(math.Floor(float64(plannedBatch) / float64(r.MaxRuns)))
		leftRuns := plannedBatch - int(r.MaxRuns) * copyRuns
		r.JobRuns = uint(copyRuns)
		if leftRuns > 0 {
			r.JobRuns++
		}

		//Manufacturing part
		for mi, m := range it.Materials {

			if m.Lvl == 1 {

				if r.MaxRuns != 999999 {

					it.Materials[mi].Quantity = helper.Round(float64(m.Quantity * float64(r.MaxRuns) * float64(1 - float64(r.ME) / 100)))

					it.Materials[mi].Quantity = it.Materials[mi].Quantity * float64(copyRuns)

					it.Materials[mi].Quantity = it.Materials[mi].Quantity + helper.Round(float64(m.Quantity * float64(leftRuns) * float64(1 - float64(r.ME) / 100)))

				} else {

					it.Materials[mi].Quantity = helper.Round(float64(m.Quantity * float64(plannedBatch) * float64(1 - float64(r.ME) / 100)))
				}

			} else {

				if m.ComponentBPO != 0 {

					it.Materials[mi].Quantity = helper.Round(float64(m.Quantity / float64(m.Batch) * float64(plannedBatch) * float64(1 - 10 / 100)))

				} else {

					it.Materials[mi].Quantity = helper.Round(float64(m.Quantity * float64(plannedBatch)))

				}

			}

			rm, err := market.GetRegionalAverage(m.MaterialID, settings.RegionID, regAvgPeriod * -1)
			if err != nil {
				return res, err
			}

			it.Materials[mi].SellPrice = rm.Sell * it.Materials[mi].Quantity
			it.Materials[mi].BuyPrice = rm.Buy * it.Materials[mi].Quantity
		}

		for mi, m := range it.Materials {
			if m.ComponentBPO != 0 && m.Lvl == 1 {

				it.Materials[mi].ProdPrice, _ = it.Materials.PrepareRecursiveRegional(m.ComponentBPO, settings.RegionID, m.LPath, regAvgPeriod, float64(plannedBatch))

			} else {

				it.Materials[mi].ProdPrice = 0

			}
		}

		r.RawProfit = r.RegionalAvg * float64(plannedBatch)

		for mi, m := range it.Materials {

			if m.Lvl == 1 {

				mm := m.SellPrice
				it.Materials[mi].BetterBSP = "s"

				min := map[string]float64{
					"s" : m.SellPrice,
					"p" : m.ProdPrice * 1.1,
					"b" : m.BuyPrice * 1.1,
				}

				for rr, e := range min {
					if e < mm && e != 0 {
						mm = e
						it.Materials[mi].BetterBSP = rr
					}
				}

				r.MaterialCost = r.MaterialCost + mm

			}

		}

		//Timing part
		tl, err := items.ProductionTimeline(it.TypeID, settings.RegionID, r.PlannedBatch, user)
		if err != nil {
			return res, err
		}

		res.Items = append(res.Items, r)

	}

	return res, nil
}
