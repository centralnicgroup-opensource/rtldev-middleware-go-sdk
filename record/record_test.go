package record

import (
	"reflect"
	"testing"
)

func TestGetData(t *testing.T) {
	data := map[string]string{
		"DOMAIN": "mydomain.com",
		"RATING": "1",
		"RNDINT": "321",
		"SUM":    "1",
	}
	rec := NewRecord(data)
	eq := reflect.DeepEqual(data, rec.GetData())
	if !eq {
		t.Error("TestGetData: Expected record data not matching.")
	}
}

func TestGetDataByKey(t *testing.T) {
	data := map[string]string{
		"DOMAIN": "mydomain.com",
		"RATING": "1",
		"RNDINT": "321",
		"SUM":    "1",
	}
	rec := NewRecord(data)
	_, err := rec.GetDataByKey("KEYNOTEXISTING")
	if err == nil {
		t.Error("TestGetDataByKey: Expected column name to not exist in record.")
	}
}
