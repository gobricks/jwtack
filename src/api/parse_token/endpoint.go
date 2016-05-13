package parse_token

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	b "github.com/gobricks/jwtack/src/backend"
)

type Request struct {
	Token string
	Key   string
}

type Response struct {
	Payload map[string]interface{} `json:"payload,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r Response) Error() error {
	return r.Err
}

func MakeEndpoint(s b.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		p, err := s.ParseToken(req.Token, req.Key)
		return Response{Payload: p, Err: err}, nil
	}
}