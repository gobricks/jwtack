package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"context"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gobricks/jwtack/src/app"
	b "github.com/gobricks/jwtack/src/backend"
	"github.com/gobricks/jwtack/src/api/create_token"
	"github.com/gobricks/jwtack/src/api/parse_token"
)

func Handler(app app.App, bs b.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(app.Logs.Error),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := mux.NewRouter()

	r.Handle(create_token.PathTemplate, kithttp.NewServer(
		create_token.MakeCreateTokenEndpoint(bs),
		create_token.DecodeCreateTokenRequest,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle(parse_token.PathTemplate, kithttp.NewServer(
		parse_token.MakeEndpoint(bs),
		parse_token.DecodeRequest,
		encodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	Error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	//case cargo.ErrUnknown:
	//	w.WriteHeader(http.StatusNotFound)
	//case ErrInvalidArgument:
	//	w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}