package backend

import (
	"github.com/gobricks/jwtack/src/app"
	"time"
)

func metricsMiddleware(metrics app.AppMetrics) ServiceMW {
	return func(next Service) Service {
		return metricsMW{metrics, next}
	}
}

type metricsMW struct {
	metrics app.AppMetrics
	Service
}

func (mw metricsMW) CreateToken(key string, payload map[string]interface{}, exp *time.Duration) (t string, err error) {
	defer func(begin time.Time) {
		mw.metrics.Access.CreateToken.Add(1)
		mw.metrics.Timers.CreateToken.CatchOverTime(time.Since(begin), time.Millisecond)
	}(time.Now())
	t, err = mw.Service.CreateToken(key, payload, exp)
	return
}

func (mw metricsMW) ParseToken(token string, key string) (payload map[string]interface{}, err error) {
	defer func(begin time.Time) {
		mw.metrics.Access.ParseToken.Add(1)
		mw.metrics.Timers.ParseToken.CatchOverTime(time.Since(begin), time.Millisecond)
	}(time.Now())
	payload, err = mw.Service.ParseToken(token, key)
	return
}