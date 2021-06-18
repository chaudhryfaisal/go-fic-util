package util

const (
	KeyHttpPr0xy               = "PR0XY"
	KeyHttpPr0xySub            = "/PR0XY/"
	KeyHttpPr0xyHeaderEndpoint = "Pr0xy-Endpoint"
)

var HeadersToDelete = []string{KeyHttpPr0xyHeaderEndpoint, "Accept-Encoding", "Host", "Referer", "Origin", "Cdn-Loop", "Cf-Connecting-Ip", "Cf-Ipcountry", "Cf-Ray", "Cf-Request-Id", "Cf-Visitor", "Connection", "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto", "X-Real-Ip"}
var EndpointHttpDebug = PropS("ENDPOINT_HTTP_DEBUG", "https://fic-debug.vercel.app")
var EndpointSelf = PropS("ENDPOINT_SELF", EndpointHttpDebug)
