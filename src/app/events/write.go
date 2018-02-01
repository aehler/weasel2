package events

import (
	"app/events/event"
	"app/events/worker"
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"errors"
	"fmt"
)

type Write struct {
	db 			func() *sqlx.DB
}

func NewWriter(db func() *sqlx.DB) worker.Actioner {

	return &Write{
		db : db,
	}

}

func (w *Write) Action(e event.Event) error {

	dbr := 0

	if entities[e.Object] == "" {
		return errors.New("BE entity not found")
	}

	if err != nil {

		fmt.Println("Couldn't fetch object")

		return err
	}

	e.EventData["user"] = e.User

	ed, errj := json.Marshal(e.EventData)
	if errj != nil {

		fmt.Println("Couldn't pack data")

		return errj
	}

	if dberr := w.db().Get(&dbr, "select * from weasel_tasks.create_new_alert($1, $2, $3, $4)",
		string(e.BusinessEvent),
		e.Object,
		e.ObjectId,
		string(ed),
	); dberr != nil {

		fmt.Println("Couldn't write business event", e, dberr)

		return dberr

	}

//	fmt.Println("alert might have been set", string(ed))

	return nil

}
