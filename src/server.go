package src

import (
	"net/http"
	app_logs "github.com/gobricks/jwtack/src/loggers"
	app_mertics "github.com/gobricks/jwtack/src/metrics"

	"github.com/gobricks/jwtack/src/api"
	"github.com/gobricks/jwtack/src/backend"
	"os/signal"
	"fmt"
	"os"
	"syscall"
	"golang.org/x/net/context"
)

const (
	defaultPort = "36801"
)

var AppLogs app_logs.AppLogs
type Config struct {
	Port string
}

func NewService() backend.Service  {
	AppLogs = app_logs.Load()
	return backend.InitService(AppLogs, app_mertics.Load())
}

func Run(cfg Config) {

	if cfg.Port == "" {
		cfg.Port = envString("PORT", defaultPort)
	}

	var (
		httpAddr = ":" + cfg.Port
		ctx = context.Background()
		mux = http.NewServeMux()
	)

	AppLogs = app_logs.Load()
	svc := backend.InitService(AppLogs, app_mertics.Load())
	mux.Handle("/api/v1/", api.Handler(ctx, svc, AppLogs.Access))

	errCh := make(chan error, 2)
	runServer(httpAddr, mux, errCh)
	waitSyscall(errCh)
	AppLogs.Error.Log("terminated", <-errCh)
}

func runServer(httpAddr string, h http.Handler, errCh chan error) {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	}))

	go func() {
		AppLogs.Info.Log("listen", "HTTP", "addr", httpAddr)
		errCh <- http.ListenAndServe(httpAddr, nil)
	}()
}

func waitSyscall(errCh chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errCh <- fmt.Errorf("%s", <-c)
	}()
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}