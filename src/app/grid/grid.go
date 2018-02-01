package grid

type grid struct {
	Columns []*Column
	Src Grider
}

func New(src Grider) *grid {

	return &grid{
		Columns : []*Column{},
		Src: src,
	}
}

func (g *grid) Column(cols ...*Column) {

	for _, c := range cols {
		g.Columns = append(g.Columns, c)
	}
}

func (g *grid) Context() map[string]interface {} {

	return map[string]interface {}{
		"columns" : g.Columns,
		"rows" : g.Src.GridRows(g.Columns),
	}

}
