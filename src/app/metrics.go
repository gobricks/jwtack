package app

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	"time"
)

type AppMetrics struct {
	Access AccessMetrics
	Timers TimeMetrics
}

type AccessMetrics struct {
	CreateToken metrics.Counter
	ParseToken  metrics.Counter
}

type TimeMetrics struct {
	OverTimeCounter metrics.Counter
	CreateToken     MethodTimeMetric
	ParseToken      MethodTimeMetric
}

type MethodTimeMetric struct {
	th       metrics.Histogram
	overtime metrics.Counter
}

func (m MethodTimeMetric)CatchOverTime(dur time.Duration, max time.Duration) {
	if dur > max {
		m.overtime.Add(1)
	}
	m.th.Observe(float64(dur))
}

func NewAppMetrics(cfg AppConfig) AppMetrics {
	appMetrics := AppMetrics{
		Access: AccessMetrics{
			CreateToken: expvar.NewCounter("access_CreateToken"),
			ParseToken: expvar.NewCounter("access_ParseToken"),
		},
		Timers: TimeMetrics{
			CreateToken: MethodTimeMetric{
				expvar.NewHistogram("duration_µs_CreateToken", 50),
				expvar.NewCounter("overtime_CreateToken"),
			},
			ParseToken: MethodTimeMetric{
				expvar.NewHistogram("duration_µs_ParseToken", 50),
				expvar.NewCounter("overtime_ParseToken"),
			},
		},
	}

	return appMetrics
}