package test

import (
	"testing"

	lr "github.com/hexonet/go-sdk/response/listresponse"
)

func TestGetList(t *testing.T) {
	raw := "[RESPONSE]\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=2\r\nPROPERTY[LIMIT][0]=2\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := lr.NewListResponse(raw)
	list := r.GetList()
	if len(list) != 2 {
		t.Error("TestGetList: unexpected amount of list entries")
	}

	noColkeys := len(r.GetColumnKeys())
	for _, row := range list {
		if noColkeys != len(row) {
			t.Error("TestGetList: unexpected amount of columns in list entry")
		}
	}
}

func TestIterator(t *testing.T) {
	raw := "[RESPONSE]\r\nPROPERTY[TOTAL][0]=100\r\nPROPERTY[FIRST][0]=0\r\nPROPERTY[LAST][0]=99\r\nPROPERTY[COUNT][0]=2\r\nPROPERTY[LIMIT][0]=2\r\nPROPERTY[CREATEDDATE][0]=2016-06-07 18:02:02\r\nPROPERTY[CREATEDDATE][1]=2008-03-18 09:37:25\r\nPROPERTY[FINALIZATIONDATE][0]=2017-06-08 18:02:02\r\nPROPERTY[FINALIZATIONDATE][1]=2017-03-19 09:37:25\r\nDESCRIPTION=Command completed successfully\r\nCODE=200\r\nQUEUETIME=0\r\nRUNTIME=0.12\r\nEOF\r\n"
	r := lr.NewListResponse(raw)
	list := r.GetList()
	if len(list) != 2 {
		t.Error("TestGetList: unexpected amount of list entries")
	}
	r.Current()
	if !r.HasNext() {
		t.Error("TestIterator: next page does not exist")
	}
	if r.HasPrevious() {
		t.Error("TestIterator: previous page does exist")
	}
	i := 0
	for r.HasNext() {
		r.Next()
		i++
	}
	if i != 1 {
		t.Error("TestIterator: expected iteration count differs")
	}
	if r.HasNext() {
		t.Error("TestIterator: next page does exist")
	}
	if !r.HasPrevious() {
		t.Error("TestIterator: previous page does not exist")
	}
	r.Rewind()
	if !r.HasNext() {
		t.Error("TestIterator: next page does not exist")
	}
	if r.HasPrevious() {
		t.Error("TestIterator: previous page does exist")
	}
	for r.HasNext() {
		r.Next()
	}
	for r.HasPrevious() {
		r.Previous()
		i--
	}
	if i != 0 {
		t.Error("TestIterator: expected iteration count differs")
	}
	if !r.HasNext() {
		t.Error("TestIterator: next page does not exist")
	}
	if r.HasPrevious() {
		t.Error("TestIterator: previous page does exist")
	}
}
