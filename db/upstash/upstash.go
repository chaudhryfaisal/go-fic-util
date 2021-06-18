package upstash

//https://console.upstash.com/pages/overall
import (
	"fmt"
	"github.com/chaudhryfaisal/go-fic-util/client"
	. "github.com/chaudhryfaisal/go-fic-util/util"
)

var Endpoint = PropS("DB_UP_STASH_ENDPOINT", EndpointHttpDebug)
var Token = PropS("DB_UP_STASH_TOKEN", "DB_UPSTASH_TOKEN")

func Get(key string) string {
	endpoint := fmt.Sprintf("%s/get/%s?_token=%s", Endpoint, key, Token)
	return result(client.HTTPRequest(&client.Req{Endpoint: endpoint, Type: &Response{}}))
}

func Set(key, val string) string {
	endpoint := fmt.Sprintf("%s/set/%s?_token=%s", Endpoint, key, Token)
	return result(client.HTTPRequest(&client.Req{Endpoint: endpoint, Type: &Response{}, Method: "POST", Payload: []byte(val)}))
}

func result(r client.Resp) string {
	return r.Body.(*Response).Result
}

type Response struct {
	Result string `json:"result"`
}
