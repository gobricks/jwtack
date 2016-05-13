package metrics

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
	th       metrics.TimeHistogram
	overtime metrics.Counter
}

func (m MethodTimeMetric)CatchOverTime(dur time.Duration, max time.Duration) {
	if dur > max {
		m.overtime.Add(1)
	}
	m.th.Observe(dur)
}

func Load() AppMetrics {
	var quantiles = []int{50, 90, 95, 99}
	appMetrics := AppMetrics{
		Access: AccessMetrics{
			CreateToken: expvar.NewCounter("access_CreateToken"),
			ParseToken: expvar.NewCounter("access_ParseToken"),
		},
		Timers: TimeMetrics{
			CreateToken: MethodTimeMetric{
				metrics.NewTimeHistogram(time.Microsecond,
					expvar.NewHistogram("duration_µs_CreateToken", 0, 10000, 3, quantiles...), ),
				expvar.NewCounter("overtime_CreateToken"),
			},
			ParseToken: MethodTimeMetric{
				metrics.NewTimeHistogram(time.Microsecond,
					expvar.NewHistogram("duration_µs_ParseToken", 0, 10000, 3, quantiles...), ),
				expvar.NewCounter("overtime_ParseToken"),
			},
		},
	}

	return appMetrics
}