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
var HTTPClient = HTTPBuildClient(PropS("PROXY", ""))

type Req struct {
	Endpoint                   string
	Method                     string
	Payload                    []byte
	Headers                    map[string]string
	ContentType                string
	Type                       interface{}
	HTTPClient                 *http.Client
	Request                    *http.Request
	KeepAcceptEncodingHeader   bool
	KeepDefaultHeadersToRemove bool
}
type Resp struct {
	Status  int
	Headers map[string][]string
	Body    interface{}
	Err     error
}

func HTTPRequest(r *Req) Resp {
	method := r.Method
	var payload []byte = nil
	if r.Payload != nil {
		payload = r.Payload
	} else if r.Request != nil {
		payload, _ = ioutil.ReadAll(r.Request.Body)
	}
	if method == "" {
		method = "GET"
		if r.Request != nil {
			method = r.Request.Method
		} else if payload != nil {
			method = "POST"
		}
	}

	req, err := http.NewRequest(method, r.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		Log.Errorf("Error creating request error=%v", err)
		return Resp{Err: err}
	}

	if r.Request != nil {
		for k, _ := range r.Request.Header {
			req.Header.Set(k, r.Request.Header.Get(k))
		}
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	if r.ContentType != "" {
		req.Header.Set("Content-Type", r.ContentType)
	}
	if !r.KeepAcceptEncodingHeader {
		req.Header.Del("Accept-Encoding")
	}
	if !r.KeepAcceptEncodingHeader {
		for _, h := range HeadersToDelete {
			req.Header.Del(h)
		}
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
					Log.Errorf("Error decoding json payload endpoint=%s error=%v payload=%s", r.Endpoint, err, string(respBytes))
					ret.Err = err
				}
			}
		} else {
			ret.Body = string(respBytes)
		}
	}
	return ret
}

func HTTPBuildClient(proxy string) *http.Client {
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
