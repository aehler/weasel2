package paginator

func NewLimiter(limit, current uint) Limiter {

	return Limiter{
		RowsLimit: limit,
		Current:   current,
	}
}

type Limiter struct {
	RowsLimit uint `json:"l"`
	Current   uint `json:"c,omitempty"`
}

func (l *Limiter) IsValid() bool {

	return l.RowsLimit > 0 && l.RowsLimit < 1000
}

func (l *Limiter) Limit() uint {

	return l.RowsLimit
}

func (l *Limiter) Offset() uint {

	if l.Current == 0 || l.Current == 1 {

		return 0
	}

	return l.RowsLimit * (l.Current - 1)
}
