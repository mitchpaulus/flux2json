package main

import (
	"strings"
	"testing"
)

// csvFromDocs is the single-table example straight from the annotated CSV docs.
const singleTableCSV = `#group,false,false,true,true,false,false,true,true,true,true
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string,string,string
#default,mean,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,host,region
,,0,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:52:00Z,15.43,mem,m,A,east
,,1,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:52:00Z,59.25,mem,m,B,east
,,2,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:52:00Z,52.62,mem,m,C,east
`

// multiTableCSV is the multi-table example from the docs.  The two tables have
// different schemas, separated by a blank line.
const multiTableCSV = `#group,false,false,true,true,false,false,true,true,true,true
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string,string,string
#default,_result,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,host,region
,,0,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:00:00Z,15.43,mem,m,A,east
,,1,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:00:00Z,59.25,mem,m,B,east
,,2,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:00:00Z,52.62,mem,m,C,east

#group,false,false,true,true,true,true,false,false,true,true
#datatype,string,long,string,string,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,string,string,string
#default,_result,,,,,,,,,
,result,table,_field,_measurement,_start,_stop,_time,_value,host,region
,,3,mem_level,m,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:00:00Z,ok,A,east
,,4,mem_level,m,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:00:00Z,info,B,east
,,5,mem_level,m,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:00:00Z,info,C,east
`

func TestParseSingleTable(t *testing.T) {
	tables, err := Parse(strings.NewReader(singleTableCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(tables))
	}

	tbl := tables[0]
	if len(tbl.Columns) != 10 {
		t.Errorf("expected 10 columns, got %d", len(tbl.Columns))
	}
	if len(tbl.Data) != 3 {
		t.Errorf("expected 3 rows, got %d", len(tbl.Data))
	}
}

func TestColumnTypes(t *testing.T) {
	tables, err := Parse(strings.NewReader(singleTableCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tbl := tables[0]

	want := map[string]struct{ typ, fmt string }{
		"result":       {"string", ""},
		"table":        {"long", ""},
		"_start":       {"datetime", "rfc3339"},
		"_stop":        {"datetime", "rfc3339"},
		"_time":        {"datetime", "rfc3339"},
		"_value":       {"double", ""},
		"_field":       {"string", ""},
		"_measurement": {"string", ""},
		"host":         {"string", ""},
		"region":       {"string", ""},
	}

	for _, col := range tbl.Columns {
		w, ok := want[col.Name]
		if !ok {
			t.Errorf("unexpected column %q", col.Name)
			continue
		}
		if col.Type != w.typ {
			t.Errorf("column %q: type = %q, want %q", col.Name, col.Type, w.typ)
		}
		if col.Format != w.fmt {
			t.Errorf("column %q: format = %q, want %q", col.Name, col.Format, w.fmt)
		}
	}
}

func TestColumnIndex(t *testing.T) {
	tables, err := Parse(strings.NewReader(singleTableCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, col := range tables[0].Columns {
		if col.Index != i {
			t.Errorf("column %q: Index = %d, want %d", col.Name, col.Index, i)
		}
	}
}

func TestValueConversion(t *testing.T) {
	tables, err := Parse(strings.NewReader(singleTableCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	row0 := tables[0].Data[0]

	// _value is double → float64
	if v, ok := row0["_value"].(float64); !ok || v != 15.43 {
		t.Errorf("_value: got %v (%T), want 15.43 float64", row0["_value"], row0["_value"])
	}

	// table is long → int64
	if v, ok := row0["table"].(int64); !ok || v != 0 {
		t.Errorf("table: got %v (%T), want int64(0)", row0["table"], row0["table"])
	}

	// _time is datetime → string (kept as-is)
	if v, ok := row0["_time"].(string); !ok || v != "2023-01-01T00:52:00Z" {
		t.Errorf("_time: got %v (%T), want string", row0["_time"], row0["_time"])
	}
}

func TestDefaultValueApplied(t *testing.T) {
	tables, err := Parse(strings.NewReader(singleTableCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The "result" column is empty in every data row; default is "mean".
	for i, row := range tables[0].Data {
		if v, ok := row["result"].(string); !ok || v != "mean" {
			t.Errorf("row %d result: got %v, want \"mean\"", i, row["result"])
		}
	}
}

func TestParseMultipleTables(t *testing.T) {
	tables, err := Parse(strings.NewReader(multiTableCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tables) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(tables))
	}
	if len(tables[0].Data) != 3 {
		t.Errorf("table 0: expected 3 rows, got %d", len(tables[0].Data))
	}
	if len(tables[1].Data) != 3 {
		t.Errorf("table 1: expected 3 rows, got %d", len(tables[1].Data))
	}
}

func TestMultipleTablesSecondSchemaHasNoBlankLine(t *testing.T) {
	// Two tables back-to-back with NO blank line between them.
	// The second #group annotation signals the new table boundary.
	input := `#datatype,string,long
#default,_result,
,result,table
,,0
,,1
#datatype,string,long
#default,_result2,
,result,table
,,2
`
	tables, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tables) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(tables))
	}
	if len(tables[0].Data) != 2 {
		t.Errorf("table 0: expected 2 rows, got %d", len(tables[0].Data))
	}
	if len(tables[1].Data) != 1 {
		t.Errorf("table 1: expected 1 row, got %d", len(tables[1].Data))
	}
	// Check that the default value from the second table is used.
	if v := tables[1].Data[0]["result"]; v != "_result2" {
		t.Errorf("table 1 row 0 result: got %v, want _result2", v)
	}
}

func TestBooleanConversion(t *testing.T) {
	input := `#datatype,boolean,boolean,boolean,boolean
#default,,,,
,a,b,c,d
,true,false,True,False
`
	tables, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	row := tables[0].Data[0]
	if row["a"] != true {
		t.Errorf("a: got %v, want true", row["a"])
	}
	if row["b"] != false {
		t.Errorf("b: got %v, want false", row["b"])
	}
	if row["c"] != true {
		t.Errorf("c: got %v, want true", row["c"])
	}
	if row["d"] != false {
		t.Errorf("d: got %v, want false", row["d"])
	}
}

func TestUnsignedLong(t *testing.T) {
	input := `#datatype,unsignedLong
#default,
,count
,18446744073709551615
`
	tables, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const maxUint64 = uint64(1<<64 - 1)
	if v, ok := tables[0].Data[0]["count"].(uint64); !ok || v != maxUint64 {
		t.Errorf("count: got %v (%T), want uint64 max", tables[0].Data[0]["count"], tables[0].Data[0]["count"])
	}
}

func TestEmptyCellBecomesNull(t *testing.T) {
	input := `#datatype,double
#default,
,value
,
`
	tables, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tables[0].Data[0]["value"] != nil {
		t.Errorf("expected nil for empty double, got %v", tables[0].Data[0]["value"])
	}
}

func TestErrorTable(t *testing.T) {
	// The annotated CSV spec defines this encoding for errors.
	input := `#datatype,string,long
,error,reference
,Failed to parse query,897
`
	tables, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(tables))
	}
	row := tables[0].Data[0]
	if row["error"] != "Failed to parse query" {
		t.Errorf("error: got %v", row["error"])
	}
	if row["reference"].(int64) != 897 {
		t.Errorf("reference: got %v", row["reference"])
	}
}
