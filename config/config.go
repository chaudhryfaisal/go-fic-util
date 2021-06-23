package config

import (
	"encoding/json"
	"fmt"
	"github.com/chaudhryfaisal/go-fic-util/db/upstash"
	. "github.com/chaudhryfaisal/go-fic-util/util"
	"reflect"
)

func Load(t interface{}) interface{} {
	key := Key(t)
	ret := reflect.New(reflect.TypeOf(t).Elem()).Interface()
	data := upstash.Get(key)
	Log.Infof("Loading data from DB key=%s data=%s", key, data)
	if len(data) > 10 {
		if err := json.Unmarshal([]byte(data), &ret); err != nil {
			Log.Errorf("Failed to load data=%s error=%v", data, err)
		}
	}
	return ret
}

func Key(t interface{}) string {
	return fmt.Sprintf("Config_%s", reflect.TypeOf(t).Elem().Name())
}

func Save(t interface{}) {
	key := Key(t)
	data := ToJsonString(t)
	resp := upstash.Set(key, data)
	Log.Infof("Saved data to DB key=%s data=%s resp=%s", key, data, resp)
}
