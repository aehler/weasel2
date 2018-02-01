package grid

type CellType string
type GridControl string

const (
	CellTypeString CellType = "string"
	CellTypeStringWithOffset CellType = "stringWithOffset"
	CellTypeInt CellType = "integer"
	CellTypeFloat CellType = "number"
	CellTypeDate CellType = "date"
	CellTypeUri CellType = "uri"
	CellTypeActions CellType = "actions"

	ControlPeriod GridControl = "grid_control_period"
)

type Column struct {
	Name string `json:"name"`
	Label string `json:"label"`
	Editable string `json:"editable"`
	Cell CellType `json:"cell"`
	Order int16   `json:"-"`
}
