package upstash

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

const (
	testKey = "TestKey"
)

func TestGetAndPut(t *testing.T) {
	val := strconv.FormatInt(time.Now().Unix(), 10)
	Get(testKey)
	Set(testKey, val)
	assert.Equal(t, val, Get(testKey))
}
