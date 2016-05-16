package backend

import (
	"github.com/gobricks/jwtack/src/app"
	"time"
)

func loggingMiddleware(logger app.AppLogs) ServiceMW {
	return func(next Service) Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logs app.AppLogs
	Service
}

func (mw logmw) CreateToken(key string, payload map[string]interface{}, exp *time.Duration) (t string, err error) {
	defer func(begin time.Time) {
		_ = mw.logs.Access.Log(
			"method", "CreateToken",
			"TokenLen", len(t),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	t, err = mw.Service.CreateToken(key, payload, exp)
	return
}

func (mw logmw) ParseToken(token string, key string) (payload map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = mw.logs.Access.Log(
			"method", "ParseToken",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	payload, err = mw.Service.ParseToken(token, key)
	return
}