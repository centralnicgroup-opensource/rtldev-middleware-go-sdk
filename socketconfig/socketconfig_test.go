package socketconfig

import (
	"strings"
	"testing"
)

func TestGetPOSTData(t *testing.T) {
	scfg := NewSocketConfig()
	if strings.Compare(scfg.GetPOSTData(), "") != 0 {
		t.Error("TestGetPOSTData: Expected postdata string should be empty.")
	}
}
