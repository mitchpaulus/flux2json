package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// tableState accumulates annotation rows, the header row, and data rows for
// one logical table before it is built into a Table.
type tableState struct {
	annotations map[string][]string // e.g. "datatype" -> ["string","long",...]
	headers     []string
	rows        [][]string
	hasHeader   bool
}

func newTableState() *tableState {
	return &tableState{annotations: make(map[string][]string)}
}

// Parse reads annotated CSV from r and returns all parsed tables.
//
// Go's encoding/csv reader skips blank lines, so we cannot rely on blank lines
// to detect table boundaries. Instead we scan line-by-line: when an annotation
// row (starting with '#') appears after at least one data row has been seen, we
// know the previous table has ended and a new one is beginning.
func Parse(r io.Reader) ([]Table, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1<<20), 1<<20) // handle long lines

	var tables []Table
	var state *tableState
	hasData := false // true once we've accumulated at least one data row

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parsing line %q: %w", line, err)
		}
		if len(fields) == 0 {
			continue
		}

		if strings.HasPrefix(fields[0], "#") {
			// ── Annotation row ────────────────────────────────────────────────
			// A '#' row after data rows marks a new table boundary.
			if hasData && state != nil {
				tbl, err := buildTable(state)
				if err != nil {
					return nil, err
				}
				tables = append(tables, tbl)
				state = nil
				hasData = false
			}

			if state == nil {
				state = newTableState()
			}

			name := strings.TrimPrefix(fields[0], "#")
			state.annotations[name] = fields[1:]

		} else {
			// ── Header or data row ────────────────────────────────────────────
			// The first column (index 0) is the annotation column, always empty
			// for header and data rows. Skip it and take the rest.
			if state == nil {
				state = newTableState()
			}

			cols := fields[1:] // drop leading annotation column

			if !state.hasHeader {
				state.headers = cols
				state.hasHeader = true
			} else {
				state.rows = append(state.rows, cols)
				hasData = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Finalize the last (or only) table.
	if state != nil && state.hasHeader {
		tbl, err := buildTable(state)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tbl)
	}

	return tables, nil
}

// parseLine parses a single CSV line using Go's standard parser so that
// quoted fields (which may contain commas or newlines) are handled correctly.
func parseLine(line string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(line))
	r.FieldsPerRecord = -1 // allow variable field counts
	r.LazyQuotes = true
	return r.Read()
}

// buildTable converts a completed tableState into a Table value.
func buildTable(s *tableState) (Table, error) {
	datatypes := s.annotations["datatype"]
	defaults := s.annotations["default"]

	columns := make([]Column, len(s.headers))
	for i, name := range s.headers {
		col := Column{Name: name, Index: i}
		if i < len(datatypes) {
			col.Type, col.Format = parseDatatype(datatypes[i])
		} else {
			col.Type = "string"
		}
		columns[i] = col
	}

	var data []map[string]any
	for _, row := range s.rows {
		record := make(map[string]any, len(columns))
		for i, col := range columns {
			raw := ""
			if i < len(row) {
				raw = row[i]
			}
			// Use the column's default value when the cell is empty.
			if raw == "" && i < len(defaults) {
				raw = defaults[i]
			}
			val, err := convertValue(raw, col)
			if err != nil {
				return Table{}, fmt.Errorf("column %q value %q: %w", col.Name, raw, err)
			}
			record[col.Name] = val
		}
		data = append(data, record)
	}

	return Table{Columns: columns, Data: data}, nil
}

// parseDatatype splits an annotated-CSV datatype string into a JSON-friendly
// type name and an optional format string.
//
//	"dateTime:RFC3339"     → ("datetime", "rfc3339")
//	"dateTime:RFC3339Nano" → ("datetime", "rfc3339nano")
//	"dateTime:number"      → ("datetime", "number")
//	"dateTime"             → ("datetime", "")
//	"double"               → ("double",   "")
func parseDatatype(dt string) (typeName, format string) {
	base, suffix, _ := strings.Cut(dt, ":")
	switch base {
	case "dateTime", "time":
		typeName = "datetime"
		format = strings.ToLower(suffix)
	default:
		typeName = base
	}
	return
}

// convertValue coerces a raw CSV string into the appropriate Go type based on
// the column descriptor.  An empty raw string becomes JSON null (nil).
func convertValue(raw string, col Column) (any, error) {
	if raw == "" {
		return nil, nil
	}

	switch col.Type {
	case "boolean":
		switch strings.ToLower(raw) {
		case "true", "t", "yes", "y", "1":
			return true, nil
		case "false", "f", "no", "n", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("unrecognized boolean value")
		}

	case "long":
		v, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, err
		}
		return v, nil

	case "unsignedLong":
		v, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, err
		}
		return v, nil

	case "double":
		v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, err
		}
		return v, nil

	default:
		// string, datetime, duration, measurement, tag, field, ignored → string
		return raw, nil
	}
}
