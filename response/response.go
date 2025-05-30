// Copyright (c) 2018 Kai Schwarz (HEXONET GmbH). All rights reserved.
//
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package response provides extended functionality to handle API response data
package response

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/column"
	"github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/record"
	rp "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responseparser"
	rt "github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/v5/responsetranslator"
)

// Response is a struct used to cover basic functionality to work with
// API response data (or hardcoded API response data).
type Response struct {
	Raw         string
	Hash        map[string]interface{}
	command     map[string]string
	columnkeys  []string
	columns     []column.Column
	recordIndex int
	records     []record.Record
}

const defaultCode = 421

// NewResponse creates a new Response object.
// It takes a raw string, a command map, and optional placeholder maps as parameters.
// The function replaces the "PASSWORD" value in the command map with "***" if it exists.
// It then translates the raw string using the command and placeholder maps.
// The function initializes a new Response object with the translated raw string, a hash value,
// the command map, empty column keys and columns, a record index of 0, and an empty records slice.
// If the hash value contains a "PROPERTY" key, the function adds columns and records to the Response object
// based on the values in the "PROPERTY" map.
// The function returns the newly created Response object.
func NewResponse(raw string, cmd map[string]string, phs ...map[string]string) *Response {
	ph := map[string]string{}
	if len(phs) > 0 {
		ph = phs[0]
	}

	newcmd := cmd
	_, exists := newcmd["PASSWORD"]
	if exists {
		newcmd["PASSWORD"] = "***"
	}
	newraw := rt.Translate(raw, cmd, ph)
	r := &Response{
		Raw:         newraw,
		Hash:        rp.Parse(newraw),
		command:     newcmd,
		columnkeys:  []string{},
		columns:     []column.Column{},
		recordIndex: 0,
		records:     []record.Record{},
	}

	h := r.GetHash()
	if p, ok := h["PROPERTY"]; ok {
		prop := p.(map[string][]string)
		colKeys := []string{}
		for key := range prop {
			colKeys = append(colKeys, key)
		}
		count := 0
		for _, c := range colKeys {
			if d, ok := prop[c]; ok {
				r.AddColumn(c, d)
				tlen := len(d)
				if tlen > count {
					count = tlen
				}
			}
		}
		for i := 0; i < count; i++ {
			d := map[string]string{}
			for _, k := range colKeys {
				col := r.GetColumn(k)
				if col != nil {
					v, err := col.GetDataByIndex(i)
					if err == nil {
						d[k] = v
					}
				}
			}
			r.AddRecord(d)
		}
	}
	return r
}

// GetCode returns the code associated with the response.
// If the response does not have a code or if the code is not a valid integer,
// it returns the default code.
//
// The code is retrieved from the response's hash, which is a map[string]interface{}.
// The code value is expected to be stored under the key "CODE" as a string.
// If the code value is not found or is not a string, the default code is returned.
//
// If an error occurs while converting the code string to an integer,
// the default code is returned and the error is logged or handled appropriately.
//
// Returns:
// - int: The code associated with the response.
func (r *Response) GetCode() int {
	h := r.GetHash()
	if h == nil {
		return defaultCode
	}

	codeStr, ok := h["CODE"].(string)
	if !ok {
		// Log or handle the error appropriately
		return defaultCode
	}

	c, err := strconv.Atoi(codeStr)
	if err != nil {
		// Log or handle the error appropriately
		return defaultCode
	}
	return c
}

// GetDescription method to return the API response description
func (r *Response) GetDescription() string {
	h := r.GetHash()
	if h == nil {
		return ""
	}

	// Assuming there is a "DESCRIPTION" key in the hash
	desc, ok := h["DESCRIPTION"].(string)
	if !ok {
		// Log or handle the error appropriately
		return ""
	}
	return desc
}

// GetPlain method to return raw API response
func (r *Response) GetPlain() string {
	return r.Raw
}

// GetQueuetime method to return API response queuetime
func (r *Response) GetQueuetime() float64 {
	h := r.GetHash()
	if val, ok := h["QUEUETIME"]; ok {
		f, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return f
		}
	}
	return 0.00
}

// GetHash method to return API response in hash format
func (r *Response) GetHash() map[string]interface{} {
	return r.Hash
}

// GetRuntime method to return API response runtime
func (r *Response) GetRuntime() float64 {
	h := r.GetHash()
	if val, ok := h["RUNTIME"]; ok {
		f, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return f
		}
	}
	return 0.00
}

// IsError method to check if API response represents an error case
func (r *Response) IsError() bool {
	c := r.GetCode()
	return (c >= 500 && c <= 599)
}

// IsSuccess method to check if API response represents a success case
func (r *Response) IsSuccess() bool {
	c := r.GetCode()
	return (c >= 200 && c <= 299)
}

// IsTmpError method to check if current API response represents a temporary error case
func (r *Response) IsTmpError() bool {
	c := r.GetCode()
	return (c >= 400 && c <= 499)
}

// IsPending method to check if current operation is returned as pending
func (r *Response) IsPending() bool {
	cmd := r.GetCommand()

	// Check if the COMMAND is AddDomain (case-insensitive)
	if command, ok := cmd["COMMAND"]; !ok || !strings.EqualFold(command, "AddDomain") {
		return false
	}

	// Retrieve the STATUS column and check if its data equals REQUESTED (case-insensitive)
	status := r.GetColumn("STATUS")
	if status != nil {
		statusData, err := status.GetDataByIndex(0)
		if err != nil {
			return false
		}
		return strings.EqualFold(statusData, "REQUESTED")
	}

	return false
}

// AddColumn method to add a Column to the column list
func (r *Response) AddColumn(key string, data []string) *Response {
	col := column.NewColumn(key, data)
	r.columns = append(r.columns, *col)
	r.columnkeys = append(r.columnkeys, key)
	return r
}

// AddRecord method to add a record to the record list
func (r *Response) AddRecord(h map[string]string) *Response {
	rec := record.NewRecord(h)
	r.records = append(r.records, *rec)
	return r
}

// GetColumn method to get column by column name
func (r *Response) GetColumn(key string) *column.Column {
	if idx, ok := r.hasColumn(key); ok {
		return &r.columns[idx]
	}
	return nil
}

// GetColumnIndex method to get data by column name and index
func (r *Response) GetColumnIndex(key string, index int) (string, error) {
	col := r.GetColumn(key)
	if col != nil {
		d, err := col.GetDataByIndex(index)
		if err == nil {
			return d, nil
		}
	}
	return "", errors.New("column Data Index does not exist")
}

// GetColumnKeys method to get the list of column names
func (r *Response) GetColumnKeys() []string {
	return r.columnkeys
}

// GetColumns method to get the list of columns
func (r *Response) GetColumns() []column.Column {
	return r.columns
}

// GetCommand method to get the underlying API command
func (r *Response) GetCommand() map[string]string {
	return r.command
}

// GetCommandPlain method to get the underlying API command in plain text
func (r *Response) GetCommandPlain() string {
	keys := make([]string, 0, len(r.command))
	for k := range r.command {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var strBuilder strings.Builder
	for i := 0; i < len(keys); i++ {
		strBuilder.WriteString(keys[i])
		strBuilder.WriteString(" = ")
		strBuilder.WriteString(r.command[keys[i]])
		strBuilder.WriteString("\n")
	}
	return strBuilder.String()
}

// GetCurrentPageNumber method to get the page number of current list query
func (r *Response) GetCurrentPageNumber() (int, error) {
	first, ferr := r.GetFirstRecordIndex()
	limit := r.GetRecordsLimitation()
	if ferr == nil && limit > 0 {
		return int(math.Floor(float64(first)/float64(limit))) + 1, nil
	}
	return 0, errors.New("could not find current page number")
}

// GetCurrentRecord method to get record of current record index
func (r *Response) GetCurrentRecord() *record.Record {
	if r.hasCurrentRecord() {
		return &r.records[r.recordIndex]
	}
	return nil
}

// GetFirstRecordIndex method to get index of first row
func (r *Response) GetFirstRecordIndex() (int, error) {
	col := r.GetColumn("FIRST")
	if col != nil {
		f, err := col.GetDataByIndex(0)
		if err == nil {
			idx, err2 := strconv.Atoi(f)
			if err2 == nil {
				return idx, nil
			}
			return 0, errors.New("could not find first record index")
		}
	}
	tlen := len(r.records)
	if tlen > 1 {
		return 0, nil
	}
	return 0, errors.New("could not find first record index")
}

// GetLastRecordIndex method to get last record index of the current list query
func (r *Response) GetLastRecordIndex() (int, error) {
	col := r.GetColumn("LAST")
	if col != nil {
		l, err := col.GetDataByIndex(0)
		if err == nil {
			idx, err2 := strconv.Atoi(l)
			if err2 == nil {
				return idx, nil
			}
			return 0, errors.New("could not find last record index")
		}
	}
	tlen := r.GetRecordsCount()
	if tlen > 0 {
		return (tlen - 1), nil
	}
	return 0, errors.New("could not find last record index")
}

// GetListHash method to get Response as List Hash including useful meta data for tables
func (r *Response) GetListHash() map[string]interface{} {
	var lh []map[string]string
	recs := r.GetRecords()
	for _, rec := range recs {
		lh = append(lh, rec.GetData())
	}
	return map[string]interface{}{
		"LIST": lh,
		"meta": map[string]interface{}{
			"columns": r.GetColumnKeys(),
			"pg":      r.GetPagination(),
		},
	}
}

// GetNextRecord method to get next record in record list
func (r *Response) GetNextRecord() *record.Record {
	if r.hasNextRecord() {
		r.recordIndex++
		return &r.records[r.recordIndex]
	}
	return nil
}

// GetNextPageNumber method to get Page Number of next list query
func (r *Response) GetNextPageNumber() (int, error) {
	cp, err := r.GetCurrentPageNumber()
	if err != nil {
		return 0, errors.New("could not find next page number")
	}
	page := cp + 1
	pages := r.GetNumberOfPages()
	if page <= pages {
		return page, nil
	}
	return pages, nil
}

// GetNumberOfPages method to get the number of pages available for this list query
func (r *Response) GetNumberOfPages() int {
	t := r.GetRecordsTotalCount()
	limit := r.GetRecordsLimitation()
	if t > 0 && limit > 0 {
		return int(math.Ceil(float64(t) / float64(limit)))
	}
	return 0
}

// GetPagination method to get pagination data; useful for table pagination
func (r *Response) GetPagination() map[string]interface{} {
	cp, err := r.GetCurrentPageNumber()
	if err != nil {
		return nil
	}
	fr, err := r.GetFirstRecordIndex()
	if err != nil {
		return nil
	}
	lr, err := r.GetLastRecordIndex()
	if err != nil {
		return nil
	}
	np, err := r.GetNextPageNumber()
	if err != nil {
		np = cp
	}
	pp, err := r.GetPreviousPageNumber()
	if err != nil {
		pp = cp
	}
	return map[string]interface{}{
		"COUNT":        r.GetRecordsCount(),
		"CURRENTPAGE":  cp,
		"FIRST":        fr,
		"LAST":         lr,
		"LIMIT":        r.GetRecordsLimitation(),
		"NEXTPAGE":     np,
		"PAGES":        r.GetNumberOfPages(),
		"PREVIOUSPAGE": pp,
		"TOTAL":        r.GetRecordsTotalCount(),
	}
}

// GetPreviousPageNumber method to get Page Number of previous list query
func (r *Response) GetPreviousPageNumber() (int, error) {
	cp, err := r.GetCurrentPageNumber()
	if err != nil {
		return 0, err
	}
	pp := cp - 1
	if pp < 1 {
		return 0, errors.New("could not find previous page number")
	}
	return pp, nil
}

// GetPreviousRecord method to get previous record in record list
func (r *Response) GetPreviousRecord() *record.Record {
	if r.hasPreviousRecord() {
		r.recordIndex--
		return &r.records[r.recordIndex]
	}
	return nil
}

// GetRecord method to get Record at given index
func (r *Response) GetRecord(idx int) *record.Record {
	if idx >= 0 && len(r.records) > idx {
		return &r.records[idx]
	}
	return nil
}

// GetRecords method to get all records
func (r *Response) GetRecords() []record.Record {
	return r.records
}

// GetRecordsCount method to get count of rows in this response
func (r *Response) GetRecordsCount() int {
	return len(r.records)
}

// GetRecordsTotalCount method to get total count of records available for the list query
func (r *Response) GetRecordsTotalCount() int {
	col := r.GetColumn("TOTAL")
	if col != nil {
		t, err := col.GetDataByIndex(0)
		if err == nil {
			c, err2 := strconv.Atoi(t)
			if err2 == nil {
				return c
			}
		}
	}
	return r.GetRecordsCount()
}

// GetRecordsLimitation method to get limit(ation) setting of the current list query
func (r *Response) GetRecordsLimitation() int {
	col := r.GetColumn("LIMIT")
	if col != nil {
		l, err := col.GetDataByIndex(0)
		if err == nil {
			if lt, err := strconv.Atoi(l); err == nil {
				return lt
			}
		}
	}
	return r.GetRecordsCount()
}

// HasNextPage method to check if this list query has a next page
func (r *Response) HasNextPage() bool {
	cp, err := r.GetCurrentPageNumber()
	if err != nil {
		return false
	}
	np := cp + 1
	return (np <= r.GetNumberOfPages())
}

// HasPreviousPage method to check if this list query has a previous page
func (r *Response) HasPreviousPage() bool {
	cp, err := r.GetCurrentPageNumber()
	if err != nil {
		return false
	}
	pp := cp - 1
	return (pp > 0)
}

// RewindRecordList method to reset index in record list back to zero
func (r *Response) RewindRecordList() *Response {
	r.recordIndex = 0
	return r
}

// hasColumn method to check if the given column exists in column list (case-insensitive)
func (r *Response) hasColumn(key string) (int, bool) {
	lowerKey := strings.ToLower(key)
	for i, k := range r.columnkeys {
		if strings.ToLower(k) == lowerKey {
			return i, true
		}
	}
	return -1, false
}

// hasCurrentRecord method to check if the record on current record index exists
func (r *Response) hasCurrentRecord() bool {
	tlen := len(r.records)
	return (tlen > 0 &&
		r.recordIndex >= 0 &&
		r.recordIndex < tlen)
}

// hasNextRecord method to check if the record list contains a next record for the
// current record index in use
func (r *Response) hasNextRecord() bool {
	next := r.recordIndex + 1
	return (r.hasCurrentRecord() && (next < len(r.records)))
}

// hasPreviousRecord method to check if the record list contains a previous record
// for the current record index in use
func (r *Response) hasPreviousRecord() bool {
	return (r.recordIndex > 0 && r.hasCurrentRecord())
}
