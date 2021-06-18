package client

import (
	"errors"
	"fmt"
	"github.com/chaudhryfaisal/go-fic-util/util"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const EndpointTodoGet = "https://jsonplaceholder.typicode.com/todos/1"
const EndpointHTTPStatTimeout = "https://httpstat.us/200?sleep=40000"

func TestHTTPRequest(t *testing.T) {
	type args struct {
		r *Req
	}
	tests := []struct {
		name     string
		args     args
		want     Resp
		contains []string
	}{
		{"GET Request", args{&Req{Endpoint: util.EndpointHttpDebug}}, Resp{Status: 200}, []string{"GET / HTTP/1.1"}},
		{"POST Request", args{&Req{Endpoint: util.EndpointHttpDebug, Method: "POST"}}, Resp{Status: 200}, []string{"POST / HTTP/1.1"}},
		{"Error", args{&Req{Endpoint: "https://fic"}}, Resp{Err: errors.New(`Get "https://fic": dial tcp: lookup fic: no such host`)}, nil},
		{"Error", args{&Req{Endpoint: "bla"}}, Resp{Err: errors.New(`Get "bla": unsupported protocol scheme ""`)}, nil},
		{"GET Headers", args{&Req{Endpoint: util.EndpointHttpDebug, Headers: map[string]string{"HAPPY": "VERY"}}}, Resp{}, []string{"Happy: VERY", "GET / HTTP/1.1"}},
		{"POST Headers", args{&Req{Endpoint: util.EndpointHttpDebug, Method: "POST", Headers: map[string]string{"HAPPY": "VERY"}}}, Resp{}, []string{"Happy: VERY", "POST / HTTP/1.1"}},
		{"Content Type", args{&Req{Endpoint: util.EndpointHttpDebug, ContentType: "HAPPY_TYPE"}}, Resp{}, []string{"Content-Type: HAPPY_TYPE"}},
		{"Map to Type", args{&Req{Endpoint: EndpointTodoGet, Type: Todo{}}}, Resp{Body: &Todo{1, 1, "delectus aut autem", false}}, nil},
		{"Map to Type Error", args{&Req{Endpoint: util.EndpointHttpDebug, Type: Todo{}}}, Resp{Err: errors.New("invalid character 'G' looking for beginning of value")}, nil},
		{"GET Request []byte{} ", args{&Req{Endpoint: util.EndpointHttpDebug, Type: []byte{}}}, Resp{Status: 200, Body: []byte{}}, nil},
		{"Error Timeout", args{&Req{Endpoint: EndpointHTTPStatTimeout}}, Resp{Err: errors.New(fmt.Sprintf(`Get "%s": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`, EndpointHTTPStatTimeout))}, nil},
		{"Content Type", args{&Req{Endpoint: "https://httpstat.us/205", ContentType: "HAPPY_TYPE"}}, Resp{}, []string{"Content-Type: HAPPY_TYPE"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTTPRequest(tt.args.r)
			if tt.want.Status > 0 {
				assert.Equal(t, tt.want.Status, got.Status)
			}
			if nil != tt.want.Err {
				assert.Equal(t, tt.want.Err.Error(), got.Err.Error())
			} else {
				assert.Nil(t, got.Err)
			}
			if nil != tt.want.Body {
				assert.NotNil(t, got.Body)
				assert.Equal(t, reflect.TypeOf(tt.want.Body).String(), reflect.TypeOf(tt.want.Body).String())
				if byteArrayStr != reflect.TypeOf(tt.want.Body).String() {
					assert.True(t, reflect.DeepEqual(got.Body, tt.want.Body), "got:%v expected:%v", got.Body, tt.want.Body)
				}
			}
			for _, s := range tt.contains {
				assert.Contains(t, got.Body, s)
			}

		})
	}
}

type Todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
