package personal

import (
	"lib/auth"
	"app/registry"
	"lib/items"
)

type PinnedBPO struct {
	TypeID uint
	ME uint
	TE uint
}

type PinnedBPOs []PinnedBPO

func ListPinned(user *auth.User) (PinnedBPOs, error) {

	res := PinnedBPOs{}

	if _, err := registry.Registry.Session.Get(user.SessionID+"_pinned"); err == nil {

		if err := registry.Registry.Session.Unmarshal(user.SessionID+"_pinned", &res); err != nil {
			return res, err
		}

	}

	return res, nil

}

func (p PinnedBPOs) Toggle(typeID uint, user *auth.User) (bool, error) {

	for _, pi := range p {

		if pi.TypeID == typeID {

			p.Remove(typeID, user)

			return false, nil

		}

	}

	err := p.Append(typeID, user)

	return true, err

}

func (p PinnedBPOs) Append(typeID uint, user *auth.User) error {

	var itl uint = 0

	if err := registry.Registry.Connect.SQLX().Get(&itl, `select
	case bpo_purchasable and coalesce(tech_lvl, 0) = 0 when true then 1 else
	coalesce(tech_lvl, 0)
	end as itl
		from emt.production_efficiency as pe
where "typeID" = $1`, typeID); err != nil {

		return err

	}

	p = append(p, PinnedBPO{
		TypeID: typeID,
		ME: items.DefaultBPOAttr[itl][0],
		TE: items.DefaultBPOAttr[itl][1],
	})

	registry.Registry.Session.Upsert(user.SessionID+"_pinned", p)

	return nil

}

func (p PinnedBPOs) Remove(typeID uint, user *auth.User) {

	for i, pi := range p {

		if pi.TypeID == typeID {

			p = append(p[:i], p[i+1:]...)

			break
		}

	}

	registry.Registry.Session.Upsert(user.SessionID+"_pinned", p)
}