package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	. "github.com/chaudhryfaisal/go-fic-util/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const byteArrayStr = "[]uint8"

var HTTPTimeout = 30 * time.Second
var HTTPClient = BuildHTTPClient(PropS("PROXY", ""))

type Req struct {
	Endpoint    string
	Method      string
	Payload     []byte
	Headers     map[string]string
	ContentType string
	Type        interface{}
	HTTPClient  *http.Client
}
type Resp struct {
	Status  int
	Headers map[string][]string
	Body    interface{}
	Err     error
}

func HTTPRequest(r *Req) Resp {
	method := r.Method
	if method == "" {
		method = "GET"
	}
	req, err := http.NewRequest(method, r.Endpoint, bytes.NewBuffer(r.Payload))
	if err != nil {
		Log.Errorf("Error creating request error=%v", err)
		return Resp{Err: err}
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	if r.ContentType != "" {
		req.Header.Set("content-type", r.ContentType)
	}
	resp, err := HTTPClient.Do(req)
	if err != nil {
		Log.Errorf("Error making request to endpoint=%s error=%v", r.Endpoint, err)
		return Resp{Err: err}
	}
	ret := Resp{Status: resp.StatusCode, Headers: resp.Header}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf("Error reading body from upstream endpoint=%s error=%v", r.Endpoint, err)
		ret.Err = err
	} else {
		if r.Type != nil {
			t := reflect.TypeOf(r.Type)
			if byteArrayStr == t.String() {
				ret.Body = respBytes
			} else {
				ret.Body = r.Type
				err = json.Unmarshal(respBytes, &ret.Body)
				if err != nil {
					Log.Errorf("Error parsing payload endpoint=%s error=%v payload=%s", r.Endpoint, err, string(respBytes))
					ret.Err = err
				}
			}
		} else {
			ret.Body = string(respBytes)
		}
	}
	return ret
}

func BuildHTTPClient(proxy string) *http.Client {
	transport := &http.Transport{}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			Log.Panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &http.Client{Transport: transport, Timeout: HTTPTimeout}
}
