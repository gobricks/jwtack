package create_token

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"fmt"
	"io/ioutil"
	"strings"
	"bytes"
	"time"
)

var PathTemplate = "/api/v1/token"

func DecodeCreateTokenRequest(r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	bodyReader := strings.NewReader(string(body))
	var reqBody struct {
		Key     string `json: "key,omitempty" `
		Payload map[string]interface{} `json: "payload,omitempty"`
		Exp     *time.Duration `json: "exp,omitempty"`
	}
	reqBody.Key = r.Header.Get("X-Jwtack-Key")
	if err := json.NewDecoder(bodyReader).Decode(&reqBody); err != nil {
		fmt.Println("decodeEncodeTokenRequest ErrorBodyRaw:", string(body))
		return nil, err
	}

	return CreateTokenRequest{reqBody.Key, reqBody.Payload, reqBody.Exp}, nil
}

func EncodeCreateTokennRequest(r *http.Request, request interface{}) (err error) {
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

func DecodeCreateTokenResponse(resp *http.Response) (interface{}, error) {
	var response CreateTokenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}