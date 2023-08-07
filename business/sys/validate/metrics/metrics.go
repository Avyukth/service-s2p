package metrics

import "expvar"

var m *Metrics

type Metrics struct {
	Goroutines *expvar.Int
	Requests   *expvar.Int
	Errors     *expvar.Int
	Panics     *expvar.Int
}

func init() {
	m = &Metrics{
		Goroutines: expvar.NewInt("goroutines"),
		Requests:   expvar.NewInt("requests"),
		Errors:     expvar.NewInt("errors"),
		Panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

const (
	key ctxKey = 1
)
