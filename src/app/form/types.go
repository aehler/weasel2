package form

const (
	Text uint = iota + 1
	TextArea
	Select
	MultipleSelect
	CheckboxGroup
	Date
	Password
	Checkbox
	Uint
	Number
	TagList
	Hidden
	Autocomplete
)

var elementType = map[string]uint {
	"text" : Text,
	"textarea" : TextArea,
	"select" : Select,
	"multipleselect" : MultipleSelect,
	"checkbox" : CheckboxGroup,
	"date" : Date,
	"password" : Password,
	"bool" : Checkbox,
	"uint" : Uint,
	"numeric" : Number,
	"taglist" : TagList,
	"hidden" : Hidden,
	"autocomplete" : Autocomplete,
}

func MapType(key string) uint {
	return elementType[key]
}

type Options []*Option

type Option struct {
	Value uint `json:"v"`
	Label string `json:"n"`
}

type Selecter interface {
	Opts() Options
}
