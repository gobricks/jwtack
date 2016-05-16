package app

import "golang.org/x/net/context"

type App struct {
	Cfg     AppConfig
	Logs    AppLogs
	Metrics AppMetrics
	Errs    AppErrors
	Ctx     context.Context
}

func NewApp() App {
	cfg := LoadConfig()
	logs := NewAppLogs(cfg)
	metrics := NewAppMetrics(cfg)
	errs := NewAppErrors(cfg)
	ctx := context.Background()

	return App{cfg, logs, metrics, errs, ctx}
}