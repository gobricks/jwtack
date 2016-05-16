package server

import (
	"net/http"
	"github.com/gobricks/jwtack/src/app"

	"github.com/gobricks/jwtack/src/api"
	"github.com/gobricks/jwtack/src/backend"
	"os/signal"
	"fmt"
	"os"
	"syscall"
)

func NewService() backend.Service {
	return backend.InitService(app.NewApp())
}

func RunServer(app app.App) {
	svc := backend.InitService(app)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", api.Handler(app, svc))

	errCh := make(chan error, 2)
	runServer(app, mux, errCh)
	waitSyscall(errCh)
	app.Logs.Error.Log("terminated", <-errCh)
}

func runServer(app app.App, h http.Handler, errCh chan error) {

	httpAddr := ":" + app.Cfg.Port

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
		app.Logs.Info.Log("listen", "HTTP", "addr", httpAddr)
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