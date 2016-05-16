package create_token

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"fmt"
	"io/ioutil"
	"strings"
	"bytes"
	"time"
)

var PathTemplate = "/api/v1/token"

func DecodeCreateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	bodyReader := strings.NewReader(string(body))
	var reqBody struct {
		Key     string `json: "key,omitempty" `
		Payload map[string]interface{} `json: "payload,omitempty"`
		Exp     *time.Duration `json: "exp_sec,omitempty"`
	}
	reqBody.Key = r.Header.Get("X-Jwtack-Key")
	if err := json.NewDecoder(bodyReader).Decode(&reqBody); err != nil {
		fmt.Println("decodeEncodeTokenRequest ErrorBodyRaw:", string(body))
		return nil, err
	}

	if reqBody.Exp != nil {
		expSec := *reqBody.Exp * time.Second
		reqBody.Exp = &expSec
	}

	return CreateTokenRequest{reqBody.Key, reqBody.Payload, reqBody.Exp}, nil
}

func EncodeCreateTokennRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	cbr := request.(CreateTokenRequest)
	if err = json.NewEncoder(&buf).Encode(cbr); err != nil {
		return err
	}

	mr := mux.NewRouter()
	mr.NewRoute().BuildOnly()
	u, err := mr.Path(PathTemplate).URLPath()
	if err != nil {
		return
	}

	r.URL.Path = u.Path
	r.Body = ioutil.NopCloser(&buf)

	fmt.Println("r.Body", r.Body)

	return
}

func DecodeCreateTokenResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response CreateTokenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}