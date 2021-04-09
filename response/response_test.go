package response

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	RP "github.com/hexonet/go-sdk/v3/responseparser"
	RTM "github.com/hexonet/go-sdk/v3/responsetemplatemanager"
)

var rtm = RTM.GetInstance()

func TestMain(m *testing.M) {
	rtm.AddTemplate(
		"login200",
		"[RESPONSE]\r\nPROPERTY[SESSION][0]=h8JLZZHdF2WgWWXlwbKWzEG3XrzoW4yshhvtqyg0LCYiX55QnhgYX9cB0W4mlpbx\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.169\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"listP0",
		"[RESPONSE]\r\nPROPERTY[TOTAL][0]=2701\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[DOMAIN][0]=0-60motorcycletimes.com\r\nPROPERTY[DOMAIN][1]=0-be-s01-0.com\r\nPROPERTY[COUNT][0]=2\r\nPROPERTY[LAST][0]=1\r\nPROPERTY[LIMIT][0]=2\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.023\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"OK",
		rtm.GenerateTemplate("200", "Command completed successfully"),
	)
	os.Exit(m.Run())
}

func TestPlaceHolderReplacements(t *testing.T) {
	r := NewResponse("", map[string]string{
		"COMMAND": "StatusAccount",
	})
	re := regexp.MustCompile(`\{[A-Z_]+\}`)
	if re.MatchString(r.GetDescription()) {
		t.Error("TestPlaceHolderReplacements: place holders not removed.")
	}

	r = NewResponse("", map[string]string{"COMMAND": "StatusAccount"}, map[string]string{"CONNECTION_URL": "123HXPHFOUND123"})
	re = regexp.MustCompile(`123HXPHFOUND123`)
	if !re.MatchString(r.GetDescription()) {
		t.Error("TestPlaceHolderReplacements: CONNECTION_URL place holder not removed.\n" + r.GetDescription())
	}
}

func TestGetCommandPlain(t *testing.T) {
	r := NewResponse("", map[string]string{
		"COMMAND": "QueryDomainOptions",
		"DOMAIN0": "example.com",
		"DOMAIN1": "example.net",
	})
	expected := "COMMAND = QueryDomainOptions\nDOMAIN0 = example.com\nDOMAIN1 = example.net\n"
	if strings.Compare(r.GetCommandPlain(), expected) != 0 {
		t.Error("TestGetCommandPlain: plain text command not matching expected value.\n\n" + r.GetCommandPlain())
	}
}

func TestGetCommandPlainSecure(t *testing.T) {
	r := NewResponse("", map[string]string{
		"COMMAND":  "CheckAuthentication",
		"SUBUSER":  "test.user",
		"PASSWORD": "test.passw0rd",
	})
	expected := "COMMAND = CheckAuthentication\nPASSWORD = ***\nSUBUSER = test.user\n"
	if strings.Compare(r.GetCommandPlain(), expected) != 0 {
		t.Error("TestGetCommandPlainSecure: plain text command not matching expected value.\n\n" + r.GetCommandPlain() + "\n\n" + expected)
	}
}

func TestGetCurrentPageNumber(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	v, err := r.GetCurrentPageNumber()
	if err != nil || v != 1 {
		t.Error(fmt.Sprintf("TestGetCurrentPageNumber: Expected current page number '%d' to be '1'.", v))
	}
}

func TestGetCurrentPageNumber2(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetCurrentPageNumber()
	if err == nil {
		t.Error("TestGetCurrentPageNumber2: Expected current page number to be error.")
	}
}

func TestGetFirstRecordIndex1(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetFirstRecordIndex()
	if err == nil {
		t.Error("TestGetFirstRecordIndex1: Expected to run into error.")
	}
}

func TestGetFirstRecordIndex2(t *testing.T) {
	h := rtm.GetTemplate("OK").GetHash()
	h["PROPERTY"] = map[string][]string{
		"DOMAIN": {"mydomain1.com", "mydomain2.com"},
	}
	serialized := RP.Serialize(h)
	r := NewResponse(serialized, map[string]string{"COMMAND": "QueryDomainList"})
	v, err := r.GetFirstRecordIndex()
	fmt.Println(serialized)
	if err != nil {
		t.Error("TestGetFirstRecordIndex2: Expected not to run into error.")
	}
	if v != 0 {
		t.Error(fmt.Printf("TestGetFirstRecordIndex2: Expected index value '%d' to be '0'.", v))
	}
}

func TestGetColumms(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	cols := r.GetColumns()
	if len(cols) != 6 {
		t.Error("Expected column size not matching.")
	}
}

func TestGetColumnIndex1(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	data, err := r.GetColumnIndex("DOMAIN", 0)
	if err != nil {
		t.Error("Expected not to run into error.")
	}
	if strings.Compare(data, "0-60motorcycletimes.com") != 0 {
		t.Error("Expected domain name not matching.")
	}
}

func TestGetColumnIndex2(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetColumnIndex("COLUMN_NOT_EXISTS", 0)
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetColumnKeys(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	colkeys := r.GetColumnKeys()
	if len(colkeys) != 6 {
		t.Error("Expected amount of columns not matching.")
	}
	defaultones := []string{"COUNT", "DOMAIN", "FIRST", "LAST", "LIMIT", "TOTAL"}
	for _, k := range defaultones {
		found := false
		for _, k2 := range colkeys {
			if strings.Compare(k, k2) == 0 {
				found = true
			}
		}
		if !found {
			t.Error(fmt.Printf("Expected column '%s' to exists.", k))
		}
	}
}

func TestGetCurrentRecord(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	rec := r.GetCurrentRecord()
	d := rec.GetData()
	expected := map[string]string{
		"COUNT":  "2",
		"DOMAIN": "0-60motorcycletimes.com",
		"FIRST":  "0",
		"LAST":   "1",
		"LIMIT":  "2",
		"TOTAL":  "2701",
	}
	eq := reflect.DeepEqual(d, expected)
	if !eq {
		t.Error("Expected data map not matching.")
	}
}

func TestGetCurrentRecord2(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	rec := r.GetCurrentRecord()
	if rec != nil {
		t.Error("Expected record to be nil.")
	}
}

func TestGetListHash(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	lh := r.GetListHash()
	if _, ok := lh["LIST"]; !ok {
		t.Error("Expected property 'LIST' to exist.")
	}
	if _, ok := lh["meta"]; !ok {
		t.Error("Expected property 'meta' to exist.")
	}
	if len(lh["LIST"].([]map[string]string)) != 2 {
		t.Error("Expected length of LIST not matching.")
	}
	cols := r.GetColumnKeys()
	mcols := lh["meta"].(map[string]interface{})["columns"].([]string)
	if !reflect.DeepEqual(mcols, cols) {
		t.Error("Expected list of columns not matching.")
	}
	pg := r.GetPagination()
	mpg := lh["meta"].(map[string]interface{})["pg"].(map[string]interface{})
	if !reflect.DeepEqual(mpg, pg) {
		t.Error("Expected pagination data not matching.")
	}
}

func TestGetNextRecord(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	rec := r.GetNextRecord()
	expected := map[string]string{"DOMAIN": "0-be-s01-0.com"}
	if !reflect.DeepEqual(rec.GetData(), expected) {
		t.Error("Expected record data not matching.")
	}
	rec = r.GetNextRecord()
	if rec != nil {
		t.Error("Expected record to be nil.")
	}
}

func TestGetPagination(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	pager := r.GetPagination()
	if pager == nil {
		t.Error("Expected pager not to be nil.")
	} else {
		pgkeys := []string{"COUNT", "CURRENTPAGE", "FIRST", "LAST", "LIMIT", "NEXTPAGE", "PAGES", "PREVIOUSPAGE", "TOTAL"}
		for _, k := range pgkeys {
			found := false
			for k2 := range pager {
				if strings.Compare(k, k2) == 0 {
					found = true
				}
			}
			if !found {
				t.Error(fmt.Printf("Expected property '%s' to exist.", k))
			}
		}
	}
}

func TestGetPreviousRecord(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	r.GetNextRecord()
	d := r.GetPreviousRecord().GetData()
	expected := map[string]string{
		"COUNT":  "2",
		"DOMAIN": "0-60motorcycletimes.com",
		"FIRST":  "0",
		"LAST":   "1",
		"LIMIT":  "2",
		"TOTAL":  "2701",
	}
	if !reflect.DeepEqual(d, expected) {
		t.Error("Expected previous record data not matching.")
	}
	d2 := r.GetPreviousRecord()
	if d2 != nil {
		t.Error("Expected previous record to be nil.")
	}
}

func TestHasNextPage1(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if r.HasNextPage() {
		t.Error("Expected no next page to exist.")
	}
}

func TestHasNextPage2(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if !r.HasNextPage() {
		t.Error("Expected next page to exist.")
	}
}

func TestHasPreviousPage1(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if r.HasPreviousPage() {
		t.Error("Expected no previous page to exist.")
	}
}

func TestHasPreviousPage2(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if r.HasPreviousPage() {
		t.Error("Expected no previous page to exist.")
	}
}

func TestGetLastRecordIndex1(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetLastRecordIndex()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetLastRecordIndex2(t *testing.T) {
	h := rtm.GetTemplate("OK").GetHash()
	h["PROPERTY"] = map[string][]string{
		"DOMAIN": {"mydomain1.com", "mydomain2.com"},
	}
	r := NewResponse(RP.Serialize(h), map[string]string{"COMMAND": "QueryDomainList"})
	lr, err := r.GetLastRecordIndex()
	if err != nil {
		t.Error("Expected not to run into error.")
	}
	if lr != 1 {
		t.Error(fmt.Printf("Expected last record index '%d' to be '1'.", lr))
	}
}

func TestGetNextPageNumber1(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetNextPageNumber()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetNextPageNumber2(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	np, err := r.GetNextPageNumber()
	if err != nil {
		t.Error("Expected not to run into error.")
	}
	if np != 2 {
		t.Error(fmt.Printf("Expected next page '%d' to be '2'.", np))
	}
}

func TestGetNumberOfPages(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	pgs := r.GetNumberOfPages()
	if pgs != 0 {
		t.Error(fmt.Printf("Expected number of pages '%d' to be '0'.", pgs))
	}
}

func TestGetPreviousPageNumber(t *testing.T) {
	plain := rtm.GetTemplate("OK").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetPreviousPageNumber()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetPreviousPageNumber2(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetPreviousPageNumber()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestRewindRecordList(t *testing.T) {
	plain := rtm.GetTemplate("listP0").GetPlain()
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	pr := r.GetPreviousRecord()
	if pr != nil {
		t.Error("Expected previous record to be nil.")
	}
	nr := r.GetNextRecord()
	if nr == nil {
		t.Error("Expected next record not to be nil.")
	}
	nr = r.GetNextRecord()
	if nr != nil {
		t.Error("Expected next record to be nil.")
	}
	pr = r.RewindRecordList().GetPreviousRecord()
	if pr != nil {
		t.Error("Expected previous record to be nil.")
	}
}
