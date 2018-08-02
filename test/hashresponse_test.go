package test

import (
	"strings"
	"testing"

	hr "github.com/hexonet/go-sdk/response/hashresponse"
)

func TestGetRaw(t *testing.T) {
	raw1 := "[RESPONSE]\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=1\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := hr.NewHashResponse(raw1)
	if strings.Compare(raw1, r.GetRaw()) != 0 {
		t.Error("TestGetRaw: value for raw response differs. [Inactive ColumnFilter]")
	}
	r.EnableColumnFilter("^CREATEDDATE$")
	plain := r.GetRaw()
	r2 := hr.NewHashResponse(plain)
	if strings.Compare(plain, r2.GetRaw()) != 0 && len(r2.GetColumnKeys()) != 1 {
		t.Error("TestGetRaw: value for raw response differs. [Active ColumnFilter]")
	}
	r.DisableColumnFilter()
	if strings.Compare(raw1, r.GetRaw()) != 0 {
		t.Error("TestGetRaw: value for raw response differs. [Reset ColumnFilter]")
	}
}

func TestGetHash(t *testing.T) {
	raw := "[RESPONSE]\r\nCODE=200\r\nDESCRIPTION=Command completed successfully\r\nRUNTIME=0.12\r\nQUEUETIME=0\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nEOF\r\n"
	r := hr.NewHashResponse(raw)
	h := r.GetHash()
	if strings.Compare(h["CODE"].(string), "200") != 0 {
		t.Error("TestGetHash: response code differs")
	}
	if strings.Compare(h["DESCRIPTION"].(string), "Command completed successfully") != 0 {
		t.Error("TestGetHash: response description differs")
	}
	if strings.Compare(h["QUEUETIME"].(string), "0") != 0 {
		t.Error("TestGetHash: response queuetime differs")
	}
	if strings.Compare(h["RUNTIME"].(string), "0.12") != 0 {
		t.Error("TestGetHash: response runtime differs")
	}
	properties := h["PROPERTY"]
	if properties == nil {
		t.Error("TestGetHash: PROPERTY key not found")
	}
	col := properties.(map[string][]string)["CREATEDDATE"]
	if col == nil {
		t.Error("TestGetHash: column key CREATEDDATE not found")
	}
	if len(col) != 2 {
		t.Error("TestGetHash: column size is other than 2")
	}
	if col[0] != "2016-06-07 18:02:02" {
		t.Error("TestGetHash: column value for index 0 differs")
	}
	if col[1] != "2008-03-18 09:37:25" {
		t.Error("TestGetHash: column value for index 1 differs")
	}
	col = properties.(map[string][]string)["FINALIZATIONDATE"]
	if col == nil {
		t.Error("TestGetHash: column key FINALIZATIONDATE not found")
	}
	if len(col) != 2 {
		t.Error("TestGetHash: column size is other than 2")
	}
	if col[0] != "2017-06-08 18:02:02" {
		t.Error("TestGetHash: column value for index 0 differs")
	}
	if col[1] != "2017-03-19 09:37:25" {
		t.Error("TestGetHash: column value for index 1 differs")
	}

	r.EnableColumnFilter("^CREATEDDATE$")
	h = r.GetHash()
	if strings.Compare(h["CODE"].(string), "200") != 0 {
		t.Error("TestGetHash: response code differs")
	}
	if strings.Compare(h["DESCRIPTION"].(string), "Command completed successfully") != 0 {
		t.Error("TestGetHash: response description differs")
	}
	if strings.Compare(h["QUEUETIME"].(string), "0") != 0 {
		t.Error("TestGetHash: response queuetime differs")
	}
	if strings.Compare(h["RUNTIME"].(string), "0.12") != 0 {
		t.Error("TestGetHash: response runtime differs")
	}
	properties = h["PROPERTY"]
	if properties == nil {
		t.Error("TestGetHash: PROPERTY key not found")
	}
	col = properties.(map[string][]string)["CREATEDDATE"]
	if col == nil {
		t.Error("TestGetHash: column key CREATEDDATE not found")
	}
	if len(col) != 2 {
		t.Error("TestGetHash: column size is other than 2")
	}
	if col[0] != "2016-06-07 18:02:02" {
		t.Error("TestGetHash: column value for index 0 differs")
	}
	if col[1] != "2008-03-18 09:37:25" {
		t.Error("TestGetHash: column value for index 1 differs")
	}
	col = properties.(map[string][]string)["FINALIZATIONDATE"]
	if col != nil {
		t.Error("TestGetHash: column key FINALIZATIONDATE should not exist")
	}

	r.DisableColumnFilter()
	h = r.GetHash()
	if strings.Compare(h["CODE"].(string), "200") != 0 {
		t.Error("TestGetHash: response code differs")
	}
	if strings.Compare(h["DESCRIPTION"].(string), "Command completed successfully") != 0 {
		t.Error("TestGetHash: response description differs")
	}
	if strings.Compare(h["QUEUETIME"].(string), "0") != 0 {
		t.Error("TestGetHash: response queuetime differs")
	}
	if strings.Compare(h["RUNTIME"].(string), "0.12") != 0 {
		t.Error("TestGetHash: response runtime differs")
	}
	properties = h["PROPERTY"]
	if properties == nil {
		t.Error("TestGetHash: PROPERTY key not found")
	}
	col = properties.(map[string][]string)["CREATEDDATE"]
	if col == nil {
		t.Error("TestGetHash: column key CREATEDDATE not found")
	}
	if len(col) != 2 {
		t.Error("TestGetHash: column size is other than 2")
	}
	if col[0] != "2016-06-07 18:02:02" {
		t.Error("TestGetHash: column value for index 0 differs")
	}
	if col[1] != "2008-03-18 09:37:25" {
		t.Error("TestGetHash: column value for index 1 differs")
	}
	col = properties.(map[string][]string)["FINALIZATIONDATE"]
	if col == nil {
		t.Error("TestGetHash: column key FINALIZATIONDATE not found")
	}
	if len(col) != 2 {
		t.Error("TestGetHash: column size is other than 2")
	}
	if col[0] != "2017-06-08 18:02:02" {
		t.Error("TestGetHash: column value for index 0 differs")
	}
	if col[1] != "2017-03-19 09:37:25" {
		t.Error("TestGetHash: column value for index 1 differs")
	}
}

func TestPaginationFunc(t *testing.T) {
	raw := "[RESPONSE]\r\nPROPERTY[LIMIT][0]=2\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=2\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := hr.NewHashResponse(raw)
	if r.Code() != 200 {
		t.Error("TestPaginationFunc: response code differs")
	}
	if r.First() != 0 {
		t.Error("TestPaginationFunc: pagination value for \"first\" differs")
	}
	if r.Total() != 100 {
		t.Error("TestPaginationFunc: pagination value for \"total\" differs")
	}
	if r.Last() != 99 {
		t.Error("TestPaginationFunc: pagination value for \"last\" differs")
	}
	if r.Limit() != 2 {
		t.Error("TestPaginationFunc: pagination value for \"limit\" differs")
	}
	if r.Count() != 2 {
		t.Error("TestPaginationFunc: pagination value for \"count\" differs")
	}
	if strings.Compare(r.Description(), "Command completed successfully") != 0 {
		t.Error("TestPaginationFunc: response description differs")
	}
	if r.Runtime() != 0.12 {
		t.Error("TestPaginationFunc: response runtime differs")
	}
	if r.Queuetime() != 0.00 {
		t.Error("TestPaginationFunc: response queuetime differs")
	}
	if r.Pages() != 50 {
		t.Error("TestPaginationFunc: pagination value for \"pages\" differs")
	}
	if r.Page() != 1 {
		t.Error("TestPaginationFunc: pagination value for \"page\" differs")
	}
	if r.Prevpage() != 1 {
		t.Error("TestPaginationFunc: pagination value for \"prevpage\" differs")
	}
	if r.Nextpage() != 2 {
		t.Error("TestPaginationFunc: pagination value for \"nextpage\" differs")
	}

	pager := r.GetPagination()
	if pager["FIRST"] != r.First() {
		t.Error("TestPaginationFunc: pagination value for \"first\" differs")
	}
	if pager["TOTAL"] != r.Total() {
		t.Error("TestPaginationFunc: pagination value for \"total\" differs")
	}
	if pager["LAST"] != r.Last() {
		t.Error("TestPaginationFunc: pagination value for \"last\" differs")
	}
	if pager["LIMIT"] != r.Limit() {
		t.Error("TestPaginationFunc: pagination value for \"limit\" differs")
	}
	if pager["COUNT"] != r.Count() {
		t.Error("TestPaginationFunc: pagination value for \"count\" differs")
	}
	if pager["PAGES"] != r.Pages() {
		t.Error("TestPaginationFunc: pagination value for \"pages\" differs")
	}
	if pager["PAGE"] != r.Page() {
		t.Error("TestPaginationFunc: pagination value for \"page\" differs")
	}
	if pager["PAGEPREV"] != r.Prevpage() {
		t.Error("TestPaginationFunc: pagination value for \"pageprev\" differs")
	}
	if pager["PAGENEXT"] != r.Nextpage() {
		t.Error("TestPaginationFunc: pagination value for \"pagenext\" differs")
	}
}

func TestCodeValidation(t *testing.T) {
	r := hr.NewHashResponse("[RESPONSE]\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n")
	if !r.IsSuccess() || r.IsTmpError() || r.IsError() {
		t.Error("TestCodeValidation: expected success case.")
	}

	r = hr.NewHashResponse("[RESPONSE]\r\nDESCRIPTION=Command completed successfully\r\nCODE=421\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n")
	if r.IsSuccess() || !r.IsTmpError() || r.IsError() {
		t.Error("TestCodeValidation: expected tmp error case.")
	}

	r = hr.NewHashResponse("[RESPONSE]\r\nDESCRIPTION=Command completed successfully\r\nCODE=500\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n")
	if r.IsSuccess() || r.IsTmpError() || !r.IsError() {
		t.Error("TestCodeValidation: expected error case.")
	}
}

func TestGetColumnKeys(t *testing.T) {
	raw := "[RESPONSE]\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=1\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := hr.NewHashResponse(raw)
	colKeys := r.GetColumnKeys()
	if len(colKeys) != 2 {
		t.Error("TestGetColumnKeys: amount of columns differs")
	}
	if colKeys[0] != "CREATEDDATE" && colKeys[0] != "FINALIZATIONDATE" {
		t.Error("TestGetColumnKeys: column name for index 0 differs")
	}
	if colKeys[1] != "FINALIZATIONDATE" && colKeys[1] != "CREATEDDATE" {
		t.Error("TestGetColumnKeys: column name for index 1 differs")
	}
}

func TestGetColumn(t *testing.T) {
	raw := "[RESPONSE]\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=1\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := hr.NewHashResponse(raw)
	col := r.GetColumn("CREATEDDATE")
	if len(col) != 2 {
		t.Error("TestGetColumn: column size differs")
	}
	if col[0] != "2016-06-07 18:02:02" {
		t.Error("TestGetColumn: column value for index 0 differs")
	}
	if col[1] != "2008-03-18 09:37:25" {
		t.Error("TestGetColumn: column value for index 1 differs")
	}
	col = r.GetColumn("IDONOTEXIST")
	if col != nil {
		t.Error("TestGetColumn: column should not exists")
	}
}

func TestGetColumnIndex(t *testing.T) {
	raw := "[RESPONSE]\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=1\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := hr.NewHashResponse(raw)
	val, err := r.GetColumnIndex("CREATEDDATE", 1)
	if err != nil || val != "2008-03-18 09:37:25" {
		t.Error("TestGetColumnIndex: column value for index 1 differs")
	}
	val, err = r.GetColumnIndex("CREATEDDATE", 2) // index n/a
	if err != nil && val != "" {
		t.Error("TestGetColumnIndex: column value for index 2 should be nil")
	}
}

// for now we leave tests for parse and serialize out as above methods won't work in case these methods break
