package create_token

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	b "github.com/gobricks/jwtack/src/backend"
	"time"
)

type CreateTokenRequest struct {
	Key     string
	Payload map[string]interface{}
	Exp     *time.Duration
}

type CreateTokenResponse struct {
	Token string `json:"token"`
	Err   error   `json:"error,omitempty"`
}

func (r CreateTokenResponse) Error() error {
	return r.Err
}

func MakeCreateTokenEndpoint(s b.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTokenRequest)
		t, err := s.CreateToken(req.Key, req.Payload, req.Exp)
		return CreateTokenResponse{Token: t, Err: err}, nil
	}
}