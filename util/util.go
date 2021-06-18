package util

import (
	"encoding/json"
	"os"
)


func PropS(key string, def string) string {
	val, _ := os.LookupEnv(key)
	if len(val) == 0 {
		val = def
	}
	if val != def {
		Log.Info("PropS: ", key, " Overridden from ENV")
	}
	return val
}

func ToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		Log.Errorf("ToJsonString obj:%v err:%v", obj, err)
		return ""
	}
	return string(b)
}
