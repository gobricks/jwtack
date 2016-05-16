package app

import (
	stdlog "log"
	"github.com/go-kit/kit/log"
	"os"
)

type AppLogs struct {
	Access, Error, Debug, Info log.Logger
}

func NewAppLogs(cfg AppConfig) AppLogs {

	var access_log log.Logger
	{
		access_log = log.NewLogfmtLogger(os.Stdout)
		access_log = log.NewContext(access_log).WithPrefix("type", "access").With("ts", log.DefaultTimestamp)
	}

	var error_log log.Logger
	{
		error_log = log.NewLogfmtLogger(os.Stderr)
		error_log = log.NewContext(error_log).WithPrefix("type", "error").With(
			"ts", log.DefaultTimestamp).With(
			"caller", log.DefaultCaller)

		stdlog.SetFlags(0)
		stdlog.SetOutput(log.NewStdlibAdapter(error_log)) // redirect anything using stdlib log to us
	}

	var debug_log log.Logger
	{
		debug_log = log.NewLogfmtLogger(os.Stdout)
		debug_log = log.NewContext(debug_log).WithPrefix("type", "debug").With(
			"ts", log.DefaultTimestamp)
	}

	var info_log log.Logger
	{
		info_log = log.NewLogfmtLogger(os.Stdout)
		info_log = log.NewContext(info_log).WithPrefix("type", "info").With(
			"ts", log.DefaultTimestamp)
	}

	return AppLogs{access_log, error_log, debug_log, info_log}
}