package items

import (
	"strconv"
)

const ActivityManufacturing uint = iota+1
const ActivityInvention uint = iota+8

const IndustryID uint = 3380
const AdvIndustryID uint = 3388

const battleshipConstructionTimeBonus uint = 408
const cruiserConstructionTimeBonus uint = 409
const frigateConstructionTimeBonus uint = 410
const industrialConstructionTimeBonus uint = 411
const blueprintmanufactureTimeBonus uint = 453
const advancedIndustrySkillIndustryJobTimeBonus uint = 1961
const manufactureTimePerLevel uint = 1982
const manufacturingTimeBonus uint = 440

const blueprintResearchTimeMultiplierBonus uint = 220
const inventionReverseEngineeringResearchSpeed uint = 1959
const attributeInventionCostMultiplier uint = 2563
const attributeInventionTimeMultiplier uint = 2564

const inventionPropabilityMultiplier uint = 1112

var (
	timeBonusAttributes []string
	invTimeBonusAttributes []string
	encMethodsIDS []uint
	DefaultBPOAttr map[uint][]uint
)

func init() {

	timeBonusAttributes = []string{
		strconv.Itoa(int(battleshipConstructionTimeBonus)),
		strconv.Itoa(int(cruiserConstructionTimeBonus)),
		strconv.Itoa(int(frigateConstructionTimeBonus)),
		strconv.Itoa(int(industrialConstructionTimeBonus)),
		strconv.Itoa(int(blueprintmanufactureTimeBonus)),
		strconv.Itoa(int(advancedIndustrySkillIndustryJobTimeBonus)),
		strconv.Itoa(int(manufactureTimePerLevel)),
		strconv.Itoa(int(manufacturingTimeBonus)),
	}

	invTimeBonusAttributes = []string{
		strconv.Itoa(int(blueprintResearchTimeMultiplierBonus)),
		strconv.Itoa(int(inventionReverseEngineeringResearchSpeed)),
		strconv.Itoa(int(attributeInventionCostMultiplier)),
		strconv.Itoa(int(attributeInventionTimeMultiplier)),
	}

	encMethodsIDS = []uint{3408, 21790, 21791, 23087, 23121}

	DefaultBPOAttr = map[uint][]uint{
		0 : {0, 0},
		1 : {10, 20},
		2 : {2, 4},
		3 : {2, 4},
	}

}
