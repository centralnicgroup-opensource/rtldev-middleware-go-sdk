package response

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	RTM "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responsetemplatemanager"
)

var rtm = RTM.GetInstance()

func TestMain(m *testing.M) {
	rtm.AddTemplate(
		"login200",
		"[RESPONSE]\r\nproperty[expiration date][0] = 2024-09-19 10:52:51\r\nproperty[sessionid][0] = bb7a884b09b9a674fb4a22211758ce87\r\ndescription = Command completed successfully\r\ncode = 200\r\nqueuetime = 0.004\r\nruntime = 0.023\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"listP0",
		"[RESPONSE]\r\nproperty[total][0] = 4\r\nproperty[first][0] = 0\r\nproperty[domain][0] = cnic-ssl-test1.com\r\nproperty[domain][1] = cnic-ssl-test2.com\r\nproperty[count][0] = 2\r\nproperty[last][0] = 1\r\nproperty[limit][0] = 2\r\ndescription = Command completed successfully\r\ncode = 200\r\nqueuetime = 0\r\nruntime = 0.007\r\nEOF\r\n",
	)
	rtm.AddTemplate(
		"pendingRegistration",
		"[RESPONSE]\r\ncode = 200\r\ndescription = Command completed successfully\r\nruntime = 0.44\r\nqueuetime = 0\r\n\r\nproperty[status][0] = REQUESTED\r\nproperty[updated date][0] = 2023-05-22 12:14:31.0\r\nproperty[zone][0] = se\r\nEOF\r\n",
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
	if !strings.Contains(r.GetDescription(), "123HXPHFOUND123") {
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
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	v, err := r.GetCurrentPageNumber()
	if err != nil || v != 1 {
		t.Errorf("TestGetCurrentPageNumber: Expected current page number '%d' to be '1'.", v)
	}
}

func TestGetCurrentPageNumber2(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetCurrentPageNumber()
	if err == nil {
		t.Error("TestGetCurrentPageNumber2: Expected current page number to be error.")
	}
}

func TestGetFirstRecordIndex1(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetFirstRecordIndex()
	if err == nil {
		t.Error("TestGetFirstRecordIndex1: Expected to run into error.")
	}
}

func TestGetColumms(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	cols := r.GetColumns()
	if len(cols) != 6 {
		t.Error("Expected column size not matching.")
	}
}

func TestGetColumnIndex1(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	data, err := r.GetColumnIndex("DOMAIN", 0)
	if err != nil {
		t.Error("Expected not to run into error.")
	}
	if strings.Compare(data, "cnic-ssl-test1.com") != 0 {
		t.Error("Expected domain name not matching.")
	}
}

func TestGetColumnIndex2(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetColumnIndex("COLUMN_NOT_EXISTS", 0)
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetColumnKeys(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
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
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	rec := r.GetCurrentRecord()
	d := rec.GetData()
	expected := map[string]string{
		"COUNT":  "2",
		"DOMAIN": "cnic-ssl-test1.com",
		"FIRST":  "0",
		"LAST":   "1",
		"LIMIT":  "2",
		"TOTAL":  "4",
	}
	eq := reflect.DeepEqual(d, expected)
	if !eq {
		t.Error("Expected data map not matching.")
	}
}

func TestGetCurrentRecord2(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	rec := r.GetCurrentRecord()
	if rec != nil {
		t.Error("Expected record to be nil.")
	}
}

func TestGetListHash(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
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
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	rec := r.GetNextRecord()
	expected := map[string]string{"DOMAIN": "cnic-ssl-test2.com"}
	if !reflect.DeepEqual(rec.GetData(), expected) {
		t.Error("Expected record data not matching.")
	}
	rec = r.GetNextRecord()
	if rec != nil {
		t.Error("Expected record to be nil.")
	}
}

func TestGetPagination(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
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
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	r.GetNextRecord()
	d := r.GetPreviousRecord().GetData()
	expected := map[string]string{
		"COUNT":  "2",
		"DOMAIN": "cnic-ssl-test1.com",
		"FIRST":  "0",
		"LAST":   "1",
		"LIMIT":  "2",
		"TOTAL":  "4",
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
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if r.HasNextPage() {
		t.Error("Expected no next page to exist.")
	}
}

func TestHasNextPage2(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if !r.HasNextPage() {
		t.Error("Expected next page to exist.")
	}
}

func TestHasPreviousPage1(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if r.HasPreviousPage() {
		t.Error("Expected no previous page to exist.")
	}
}

func TestHasPreviousPage2(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	if r.HasPreviousPage() {
		t.Error("Expected no previous page to exist.")
	}
}

func TestGetLastRecordIndex1(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetLastRecordIndex()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetNextPageNumber1(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetNextPageNumber()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetNextPageNumber2(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
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
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	pgs := r.GetNumberOfPages()
	if pgs != 0 {
		t.Error(fmt.Printf("Expected number of pages '%d' to be '0'.", pgs))
	}
}

func TestGetPreviousPageNumber(t *testing.T) {
	plain := rtm.GetTemplate("OK")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetPreviousPageNumber()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestGetPreviousPageNumber2(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
	r := NewResponse(plain, map[string]string{"COMMAND": "QueryDomainList"})
	_, err := r.GetPreviousPageNumber()
	if err == nil {
		t.Error("Expected to run into error.")
	}
}

func TestRewindRecordList(t *testing.T) {
	plain := rtm.GetTemplate("listP0")
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

func TestIsPending(t *testing.T) {
	plain := rtm.GetTemplate("pendingRegistration")
	r := NewResponse(plain, map[string]string{"COMMAND": "AddDomain"})
	fmt.Print(r.IsPending())
	if got := r.IsPending(); got != true {
		t.Errorf("isPending() = %v, want true", got)
	}
}
