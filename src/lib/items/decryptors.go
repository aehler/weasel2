package items

import (
	"app/registry"
	"sync"
	"fmt"
)

type Decriptor struct {
	TypeID uint `db:"typeID"`
	TypeName string `db:"typeName"`
	InvChance float64 `db:"chance"`
	MaxRuns int `db:"max_run"`
	ME int `db:"me"`
	TE int `db:"te"`
	AdjPrice float64 `db:"adj_price"`
}

var decryptors struct {
	d []Decriptor
	mu sync.Mutex
}

func GetDecryptorData() error {

	decryptors.mu.Lock()
	defer decryptors.mu.Unlock()

	if err := registry.Registry.Connect.SQLX().Select(&decryptors.d, `with tt as (select t."typeID", t."typeName",
case when ta."attributeID" = 1112 then ta."valueFloat" else 0 end as chance,
case when ta."attributeID" = 1113 then ta."valueFloat" else 0 end as me,
case when ta."attributeID" = 1114 then ta."valueFloat" else 0 end as te,
case when ta."attributeID" = 1124 then ta."valueFloat" else 0 end as max_run
 from evesde."invTypes" as t
left join evesde."dgmTypeAttributes" as ta using("typeID")
where t."groupID" = 1304)

select "typeID", "typeName", sum(chance) as chance, sum(me) as me, sum(te) as te, sum(max_run) as max_run,
ma.adj_price
from tt
left join emt.market_avg as ma on "typeID" = type_id
group by "typeID", "typeName", ma.adj_price`); err != nil {

		return err

	}

	fmt.Println("Got decryptor data")

	return nil
}

func Decryptors() []Decriptor {

	decryptors.mu.Lock()
	defer decryptors.mu.Unlock()

	return decryptors.d
}