package http

import (
	. "chaudhryfaisal/go-fic-util/com"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	KeyHttpPr0xy               = "PR0XY"
	KeyHttpPr0xySub            = "/PR0XY/"
	KeyHttpPr0xyHeaderEndpoint = "PR0XY-Endpoint"
)

var RegexUrl = regexp.MustCompile(`(https?)://([\w.-]+):?(\d+)?([/.a-zA-Z0-9?=&_~:#\[\]@!$'()*+;%{}-]+)?`)
var HTTPClient = BuildHTTPClient(PropS("PROXY", ""))

func HttpRespondJson(rw http.ResponseWriter, json string) {
	HttpRespondJsonWithStatus(rw, json, http.StatusOK)
}
func HttpRespondJsonWithStatus(rw http.ResponseWriter, json string, status int) {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write([]byte(json))
}
func HTTPProxyPass(rw http.ResponseWriter, r *http.Request, endpoint string) {
	req, err := http.NewRequest(r.Method, endpoint, r.Body)
	HTTPCopyHeadersForProxy(r, req)
	resp, err := HTTPClient.Do(req)
	if err != nil {
		HTTPError(err, rw, "HttpProxy", "Error proxying to upstream", "endpoint="+endpoint)
		return
	}
	rw.WriteHeader(resp.StatusCode)
	HTTPCopyHeadersForClientResponse(resp, rw)
	_, _ = io.Copy(rw, resp.Body)
	_ = resp.Body.Close()
}

func HTTPCopyHeadersForClientResponse(src *http.Response, dest http.ResponseWriter) {
	for name, values := range src.Header {
		for _, value := range values {
			if strings.ToLower(name) != "content-length" {
				dest.Header().Set(name, value)
			}
		}
	}
}

func HTTPCopyHeadersForProxy(src *http.Request, dest *http.Request) {
	for name, values := range src.Header {
		for _, value := range values {
			dest.Header.Set(name, value)
		}
	}
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
	return &http.Client{Transport: transport}
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

func HTTPProxifyURLsInData(selfUrl string, data string) string {
	return HTTPProxifyURLsInDataWithFilters(selfUrl, data, nil, nil)
}
func HTTPProxifyURLsInDataWithFilters(selfUrl string, data string, include *regexp.Regexp, exclude *regexp.Regexp) string {
	urls := RegexUrl.FindAllStringSubmatch(data, -1)
	if urls == nil {
		Log.Debug("HTTPProxifyURLs: No urls found in ", data)
		return data
	}
	for _, u := range urls {
		if strings.Contains(u[0], KeyHttpPr0xySub) {
			continue
		}
		if include != nil && !include.MatchString(u[0]) {
			continue
		}
		if exclude != nil && exclude.MatchString(u[0]) {
			continue
		}
		port := u[3]
		if (len(port)) == 0 {
			if "https" == u[1] {
				port = "443"
			} else {
				port = "80"
			}
		}
		proxied := fmt.Sprintf("%s/%s/%s/%s/%s%s", selfUrl, KeyHttpPr0xy, u[1], u[2], port, u[4])
		data = strings.ReplaceAll(data, u[0], proxied)
	}
	return data
}
