package parse_token

import (
	//"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"fmt"
)

var PathTemplate = "/api/v1/token/{token}"

func DecodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		return nil, fmt.Errorf("Incorrect token %s", token)
	}

	key := r.Header.Get("X-Jwtack-Key")
	if key == "" {
		return nil, fmt.Errorf("Incorrect X-JWTAck-Key")
	}

	return Request{Token: token, Key: key}, nil
}
//
//func EncodeParseTokenRequest(r *http.Request, request interface{}) (err error) {
//	var buf bytes.Buffer
//	cbr := request.(CreateTokenRequest)
//	if err = json.NewEncoder(&buf).Encode(cbr); err != nil {
//		return err
//	}
//
//	mr := mux.NewRouter()
//	mr.NewRoute().BuildOnly()
//	u, err := mr.Path(pathTpls.createToken).URLPath()
//	if err != nil {
//		return
//	}
//
//	r.URL.Path = u.Path
//	r.Body = ioutil.NopCloser(&buf)
//
//	fmt.Println("r.Body", r.Body)
//
//	return
//}
//
//func DecodeParseTokenResponse(resp *http.Response) (interface{}, error) {
//	var response CreateTokenResponse
//	err := json.NewDecoder(resp.Body).Decode(&response)
//	return response, err
//}