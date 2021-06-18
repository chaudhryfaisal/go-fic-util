package proxy

import (
	"fmt"
	. "github.com/chaudhryfaisal/go-fic-util/client"
	. "github.com/chaudhryfaisal/go-fic-util/util"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var RegexUrl = regexp.MustCompile(`(https?)://([\w.-]+):?(\d+)?([/.a-zA-Z0-9?=&_~:#\[\]@!$'()*+;%{}-]+)?`)

func HTTPRespondJson(rw http.ResponseWriter, json string) {
	HTTPRespondJsonWithStatus(rw, json, http.StatusOK)
}
func HTTPRespondJsonWithStatus(rw http.ResponseWriter, json string, status int) {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write([]byte(json))
}
func HTTPProxyPass(rw http.ResponseWriter, r *http.Request, endpoint string) {
	req, err := http.NewRequest(r.Method, endpoint, r.Body)
	HTTPCopyHeadersForUpstream(r, req)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		HTTPError(err, rw, "HTTPProxyPass", "HTTPError proxying to upstream", "endpoint="+endpoint)
	} else {
		rw.WriteHeader(resp.StatusCode)
		HTTPCopyHeadersForClient(resp, rw)
		_, _ = io.Copy(rw, resp.Body)
		_ = resp.Body.Close()
	}
}

func HTTPCopyHeadersForClient(src *http.Response, dest http.ResponseWriter) {
	for name, values := range src.Header {
		for _, value := range values {
			if strings.ToLower(name) != "content-length" {
				dest.Header().Set(name, value)
			}
		}
	}
}
func HTTPCopyHeadersForClientFromResp(src *Resp, dest http.ResponseWriter) {
	for name, values := range src.Headers {
		for _, value := range values {
			if strings.ToLower(name) != "content-length" {
				dest.Header().Set(name, value)
			}
		}
	}
}

func HTTPCopyHeadersForUpstream(src *http.Request, dest *http.Request) {
	for name, values := range src.Header {
		for _, value := range values {
			dest.Header.Set(name, value)
		}
	}
	for _, h := range HeadersToDelete {
		dest.Header.Del(h)
	}
}

func HTTPError(err error, rw http.ResponseWriter, source, msg, details string) {
	Log.Errorf("%s: msg=%s details=%s error=%v", source, msg, details, err)
	http.Error(rw, fmt.Sprintf("{\"error:\":\"%s\"}", msg), http.StatusBadRequest)
}

func HTTPRequestPathWithQuery(r *http.Request) string {
	path := r.RequestURI
	if len(path) == 0 {
		path = r.URL.Path
	}
	if len(r.URL.RawQuery) > 0 {
		path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
	}
	return path
}
