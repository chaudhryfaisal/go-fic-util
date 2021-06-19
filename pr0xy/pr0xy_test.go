package pr0xy

import (
	"fmt"
	"github.com/chaudhryfaisal/go-fic-util/proxy"
	. "github.com/chaudhryfaisal/go-fic-util/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func Test_PR0XY(t *testing.T) {
	testPr0xyPath(t, "")
	testPr0xyPath(t, "/")
}

func testPr0xyPath(t *testing.T, path string) {
	endpoint := fmt.Sprintf("%s%s/443%s", KeyHttpPr0xySub, strings.Replace(EndpointHttpDebug, "://", "/", 1), path)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	pr0xyHandler := func(rw http.ResponseWriter, r *http.Request, endpoint string) {
		t.Log("pr0xyHandler: ", endpoint)
		proxy.HTTPProxyPass(rw, r, endpoint)
	}
	httpHandler := func(rw http.ResponseWriter, r *http.Request) {
		Pr0xy(rw, r, pr0xyHandler)
	}
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(httpHandler)
	handler.ServeHTTP(resp, req)
	// Check the status code is what we expect.
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("HTTPError reading response body ", err)
		return
	}
	body := string(respBytes)
	if path == "" {
		path = "/"
	}
	assert.Contains(t, body, "Host: fic-debug.vercel.app")
	assert.Contains(t, body, fmt.Sprintf("GET %s HTTP/1.1", path))
}

var self = "http://self.com"

func Test_HTTPPr0xifyURLsInData(t *testing.T) {
	assertEqual(t, Pr0xifyURLsInData(self, "http://google.com/"+KeyHttpPr0xy+"/"), "http://google.com/"+KeyHttpPr0xy+"/")
	assertEqual(t, Pr0xifyURLsInData(self, "http://google.com"), fmt.Sprintf("%s/PR0XY/http/google.com/80", self))
	assertEqual(t, Pr0xifyURLsInData(self, "http://google.com/"), fmt.Sprintf("%s/PR0XY/http/google.com/80/", self))
	assertEqual(t, Pr0xifyURLsInData(self, "http://google.com/bla"), fmt.Sprintf("%s/PR0XY/http/google.com/80/bla", self))
	assertEqual(t, Pr0xifyURLsInData(self, "https://google.com/bla"), fmt.Sprintf("%s/PR0XY/https/google.com/443/bla", self))
	assertEqual(t, Pr0xifyURLsInData(self, "https://google.com:8080/bla"), fmt.Sprintf("%s/PR0XY/https/google.com/8080/bla", self))
	assertEqual(t, Pr0xifyURLsInData(self, "https://global.edge.bamgrid.com/accounts/me/profiles/{profileId}"), "http://self.com/PR0XY/https/global.edge.bamgrid.com/443/accounts/me/profiles/{profileId}")
	assertEqual(t, Pr0xifyURLsInData(self, "[\"https://google.com:8080/bla\",\"https://google.com/bla\"]"), "[\"http://self.com/PR0XY/https/google.com/8080/bla\",\"http://self.com/PR0XY/https/google.com/443/bla\"]")
}
func Test_HTTPPr0xifyURLsInDataWithFilters(t *testing.T) {

	include := regexp.MustCompile(`/media/`)
	exclude := regexp.MustCompile(`/dd5f96ab-247f-4a78-8530-da019cb73d2d/`)
	assert.Equal(t, "http://self.com/PR0XY/https/hotstar.playback.edge.bamgrid.com/443/media/31789ef0-f553-4b17-93df-a622444933c0/scenarios/{scenario},https://vod-bgc-na-east-1.media.dssott.com/bgui/ps01/hotstar/bgui/2021/02/25/1614295564-326104.mp4,https://prod-ripcut-delivery.hotstar-plus.net/v1/variant/hotstar/DD981E2E05C81E71D4A3AB8974537639831CF0DF4AA094263EA7DBDB1171385B,http://self.com/PR0XY/https/hotstar.playback.edge.bamgrid.com/443/media/dd5f96ab-247f-4a78-8530-da019cb73d2d/scenarios/{scenario},http://self.com/PR0XY/https/hotstar.playback.edge.bamgrid.com/443/media/9533008d-3ed9-45c4-8077-ba6ee68aed0c/scenarios/{scenario}", Pr0xifyURLsInDataWithFilters(self, "https://hotstar.playback.edge.bamgrid.com/media/31789ef0-f553-4b17-93df-a622444933c0/scenarios/{scenario},https://vod-bgc-na-east-1.media.dssott.com/bgui/ps01/hotstar/bgui/2021/02/25/1614295564-326104.mp4,https://prod-ripcut-delivery.hotstar-plus.net/v1/variant/hotstar/DD981E2E05C81E71D4A3AB8974537639831CF0DF4AA094263EA7DBDB1171385B,https://hotstar.playback.edge.bamgrid.com/media/dd5f96ab-247f-4a78-8530-da019cb73d2d/scenarios/{scenario},https://hotstar.playback.edge.bamgrid.com/media/9533008d-3ed9-45c4-8077-ba6ee68aed0c/scenarios/{scenario}", include, nil))
	assert.Equal(t, "http://self.com/PR0XY/https/hotstar.playback.edge.bamgrid.com/443/media/31789ef0-f553-4b17-93df-a622444933c0/scenarios/{scenario},https://vod-bgc-na-east-1.media.dssott.com/bgui/ps01/hotstar/bgui/2021/02/25/1614295564-326104.mp4,https://prod-ripcut-delivery.hotstar-plus.net/v1/variant/hotstar/DD981E2E05C81E71D4A3AB8974537639831CF0DF4AA094263EA7DBDB1171385B,https://hotstar.playback.edge.bamgrid.com/media/dd5f96ab-247f-4a78-8530-da019cb73d2d/scenarios/{scenario},http://self.com/PR0XY/https/hotstar.playback.edge.bamgrid.com/443/media/9533008d-3ed9-45c4-8077-ba6ee68aed0c/scenarios/{scenario}", Pr0xifyURLsInDataWithFilters(self, "https://hotstar.playback.edge.bamgrid.com/media/31789ef0-f553-4b17-93df-a622444933c0/scenarios/{scenario},https://vod-bgc-na-east-1.media.dssott.com/bgui/ps01/hotstar/bgui/2021/02/25/1614295564-326104.mp4,https://prod-ripcut-delivery.hotstar-plus.net/v1/variant/hotstar/DD981E2E05C81E71D4A3AB8974537639831CF0DF4AA094263EA7DBDB1171385B,https://hotstar.playback.edge.bamgrid.com/media/dd5f96ab-247f-4a78-8530-da019cb73d2d/scenarios/{scenario},https://hotstar.playback.edge.bamgrid.com/media/9533008d-3ed9-45c4-8077-ba6ee68aed0c/scenarios/{scenario}", include, exclude))
}

func assertEqual(t *testing.T, actual string, expected string) {
	assert.Equal(t, expected, actual)
}
