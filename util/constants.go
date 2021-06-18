package util

const (
	KeyHttpPr0xy               = "PR0XY"
	KeyHttpPr0xySub            = "/PR0XY/"
	KeyHttpPr0xyHeaderEndpoint = "Pr0xy-Endpoint"
)

var EndpointHttpDebug = PropS("ENDPOINT_HTTP_DEBUG", "https://fic-debug.vercel.app")
var EndpointSelf = PropS("ENDPOINT_SELF", EndpointHttpDebug)
