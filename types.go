package main

// Column describes a table column's metadata.
type Column struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
	Index  int    `json:"index"`
}

// Table holds a single parsed table with its column metadata and data rows.
type Table struct {
	Columns []Column         `json:"columns"`
	Data    []map[string]any `json:"data"`
}
