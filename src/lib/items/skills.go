package items

import (
	"app/registry"
	"strings"
	"fmt"
	"errors"
	"encoding/json"
)

type Skills []Skill

type Skill struct {
	SkillID uint `json:"skillID" db:"skill_id"`
	Level uint `json:"lvl" db:"req_value"`
	SkillName string `json:"name" db:"typeName"`
	Tree Skills `json:"tree" db:"-"`
	TreeLvl uint `json:"treelvl" db:"lvl"`
	ReqID uint `json:"req_id" db:"req_id"`
	TimeBonus int `json:"time_bonus" db:"time_bonus"`
	ProbBonus int `json:"prob_bonus,omitempty" db:"prob_bonus"`
}

func (s *Skill) SubSkills(attrlist []string) error {

	res := []Skill{}

	if err := registry.Registry.Connect.SQLX().Select(&res, `with recursive r as (
	SELECT
	skill."typeID" as skill_id, skill."valueInt" as req_id, types."typeName",
	case attr."attributeID" when 182 then
		(select  coalesce(skill2."valueInt", skill2."valueFloat")::bigint as req_id FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."typeID" and skill2."attributeID" = 277)
		when 183 then
		(select  coalesce(skill2."valueInt", skill2."valueFloat")::bigint as req_id FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."typeID" and skill2."attributeID" = 278)
		when 184 then
		(select  coalesce(skill2."valueInt", skill2."valueFloat")::bigint as req_id FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."typeID" and skill2."attributeID" = 279)
	end as req_value,
	1 as lvl,
		coalesce((select coalesce(skill2."valueInt", skill2."valueFloat") as time_bonus FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."valueInt"
		and skill2."attributeID" = any($2)), 0) as time_bonus,
		coalesce((select coalesce(skill2."valueInt", skill2."valueFloat") as time_bonus FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."valueInt"
		and skill2."attributeID" = $3), 0) as prob_bonus
	FROM evesde."dgmTypeAttributes" AS skill
	LEFT JOIN evesde."dgmAttributeTypes" AS attr ON skill."attributeID" = attr."attributeID"
	LEFT JOIN evesde."invTypes" AS types ON skill."valueInt" = types."typeID"
	WHERE skill."typeID" = $1
	and attr."attributeID" in (182, 183, 184)

union

	SELECT
	skill."typeID" as skill_id, skill."valueInt" as req_id, types."typeName",
	case attr."attributeID" when 182 then
		(select  coalesce(skill2."valueInt", skill2."valueFloat")::bigint as req_id FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."typeID" and skill2."attributeID" = 277)
		when 183 then
		(select  coalesce(skill2."valueInt", skill2."valueFloat")::bigint as req_id FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."typeID" and skill2."attributeID" = 278)
		when 184 then
		(select  coalesce(skill2."valueInt", skill2."valueFloat")::bigint as req_id FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."typeID" and skill2."attributeID" = 279)
	end as req_value,
	lvl+1 as lvl,
		coalesce((select coalesce(skill2."valueInt", skill2."valueFloat") as time_bonus FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."valueInt"
		and skill2."attributeID" = any($2)), 0) as time_bonus,
		coalesce((select coalesce(skill2."valueInt", skill2."valueFloat") as time_bonus FROM evesde."dgmTypeAttributes" AS skill2 where skill2."typeID" = skill."valueInt"
		and skill2."attributeID" = $3), 0) as prob_bonus
	FROM evesde."dgmTypeAttributes" AS skill
	LEFT JOIN evesde."dgmAttributeTypes" AS attr ON skill."attributeID" = attr."attributeID"
	LEFT JOIN evesde."invTypes" AS types ON skill."valueInt" = types."typeID"
	JOIN r on skill."typeID" = r.req_id
	WHERE attr."attributeID" in (182, 183, 184)
)
select * from r where req_id is not null`, s.SkillID, fmt.Sprintf("{%s}", strings.Join(attrlist, ",")), inventionPropabilityMultiplier); err != nil {

		fmt.Println("Error, skill,", err.Error(), s.SkillID, strings.Join(attrlist, ","))

		return err
	}

	s.Tree = res

	return nil

}

func (u *Skills) Scan(src interface{}) error {

	var source []byte

	switch src.(type) {

	case string:

		source = []byte(src.(string))

	case []byte:

		source = src.([]byte)

	default:

		return errors.New("Incompatible type for Skills")
	}

	if err := json.Unmarshal(source, &u); err != nil {

		return err
	}

	return nil
}