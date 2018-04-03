package response

import (
	"time"
	"errors"
	"reflect"
)

type EsiDate struct {
	time.Time
}

func (self *EsiDate) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1:len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
	}
	self.Time = t
	return
}

func (s *EsiDate) Scan(src interface{}) error {

	var source string

	switch src.(type) {

	case string:

		source = src.(string)

	case []byte:

		source = string(src.([]byte))

	case []int8:

		source = string(src.([]byte))

	case time.Time:

		*s = EsiDate{src.(time.Time)}

		return nil

	default:

		return errors.New("Incompatible type for EsiDate - "+reflect.TypeOf(src).String())

	}

	if t, err := time.Parse("2006-01-02", source); err == nil {

		*s = EsiDate{t}

	} else {

		return err

	}

	return nil
}