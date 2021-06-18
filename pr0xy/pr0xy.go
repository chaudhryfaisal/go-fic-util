package disney

import (
	"fmt"
	"github.com/chaudhryfaisal/go-fic-util/client"
	. "github.com/chaudhryfaisal/go-fic-util/proxy"
	. "github.com/chaudhryfaisal/go-fic-util/util"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var RegexPr0xy = regexp.MustCompile(`/PR0XY/(https?)/([\w-_.]+)/(\d+)?([\w+&@#/%?=~_|!:,.;-]+)?`)

type Pr0xyHandler func(rw http.ResponseWriter, r *http.Request, endpoint string)

func Pr0xy(rw http.ResponseWriter, r *http.Request, handler Pr0xyHandler) {
	path := RequestPathWithQuery(r)
	server := RegexPr0xy.FindStringSubmatch(path)
	if server == nil {
		Error(nil, rw, "Pr0xy", "Invalid Request", fmt.Sprintf(`failed to match %s to %s`, path, RegexPr0xy.String()))
		return
	}
	endpoint := fmt.Sprintf("%s://%s:%s%s", server[1], server[2], server[3], server[4])
	r.URL.Path = server[4]
	r.RequestURI = server[4]
	r.Header.Set(KeyHttpPr0xyHeaderEndpoint, endpoint)
	handler(rw, r, endpoint)
}

func Pr0xyWithReplaceInBody(rw http.ResponseWriter, r *http.Request, include *regexp.Regexp, exclude *regexp.Regexp) {
	endpoint := r.Header.Get(KeyHttpPr0xyHeaderEndpoint)
	r.Header.Del(KeyHttpPr0xyHeaderEndpoint)
	req, err := http.NewRequest(r.Method, endpoint, r.Body)
	HTTPCopyHeadersForProxy(r, req)
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		Error(err, rw, "Pr0xyWithReplaceInBody", "Error proxying to upstream", "endpoint="+endpoint)
		return
	}
	rw.WriteHeader(resp.StatusCode)
	HTTPCopyHeadersForClientResponse(resp, rw)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error(err, rw, "Pr0xyWithReplaceInBody", "error parsing payload from upstream", "respBytes="+string(respBytes))
		return
	}
	_, _ = rw.Write([]byte(Pr0xifyURLsInDataWithFilters(EndpointSelf, string(respBytes), include, exclude)))
}

func Pr0xifyURLsInData(selfUrl string, data string) string {
	return Pr0xifyURLsInDataWithFilters(selfUrl, data, nil, nil)
}

func Pr0xifyURLsInDataWithFilters(selfUrl string, data string, include *regexp.Regexp, exclude *regexp.Regexp) string {
	urls := RegexUrl.FindAllStringSubmatch(data, -1)
	if urls == nil {
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
