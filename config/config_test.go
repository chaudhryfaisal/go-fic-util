package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_load(t *testing.T) {
	test := TestStruct{Enabled: true, MultiplyPass: false, MultiplyPassCount: time.Now().Nanosecond()}
	Save(&test)
	loaded := Load(&TestStruct{})
	assert.Equal(t, &test, loaded)
}
func TestKey(t *testing.T) {
	assert.Equal(t, "Config_TestStruct", Key(&TestStruct{}))
}

type TestStruct struct {
	Enabled           bool `form:"type=checkbox"`
	MultiplyPass      bool `form:"type=checkbox"`
	MultiplyPassCount int
	Webhook           string `json:"-" form:"-" link:"http://bla.com"`
}
