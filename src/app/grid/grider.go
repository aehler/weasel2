package grid

type Grider interface{
	GridRows([]*Column) []map[string]interface {}
}
