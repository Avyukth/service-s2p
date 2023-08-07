package metrics

import (
	"context"
	"expvar"
)

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

//

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

func AddGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*Metrics); ok {
		if v.Goroutines.Value()%100 == 0 {
			v.Goroutines.Add(1)
		}
	}
}

func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*Metrics); ok {
		v.Requests.Add(1)
	}
}

//

func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*Metrics); ok {
		v.Errors.Add(1)
	}
}

//

func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*Metrics); ok {
		v.Panics.Add(1)
	}
}
