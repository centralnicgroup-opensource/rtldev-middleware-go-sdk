package column

import (
	"strings"
	"testing"
)

func TestGetKey(t *testing.T) {
	data := []string{"mydomain1.com", "mydomain2.com", "mydomain3.com"}
	col := NewColumn("DOMAIN", data)
	if strings.Compare(col.GetKey(), "DOMAIN") != 0 {
		t.Error("TestGetKey: Expected column name not matching.")
	}
}
