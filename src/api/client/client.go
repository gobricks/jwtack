package apiclient

import (
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"

	b "github.com/gobricks/jwtack/src/backend"
	"github.com/gobricks/jwtack/src/api/create_token"

	"strings"
)

func New(ctx context.Context, instance string, logger log.Logger, c *http.Client) b.Service {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	url, err := url.Parse(instance); if err != nil {
		panic(err)
	}
	return client{
		Context: ctx,
		Logger:  logger,
		createToken: httptransport.NewClient("POST", url,
			create_token.EncodeCreateTokennRequest,
			create_token.DecodeCreateTokenResponse,
			httptransport.SetClient(c),
		).Endpoint(),
	}
}

type client struct {
	context.Context
	log.Logger
	createToken    endpoint.Endpoint
}

func (c client) CreateToken(key string, payload map[string]interface{}) (t string, err error) {
	response, err := c.createToken(c.Context, create_token.CreateTokenRequest{Key: key, payload: payload})
	if err != nil {
		return
	}
	t = response.(create_token.CreateTokenResponse).Token
	err = response.(create_token.CreateTokenResponse).Err
	return
}