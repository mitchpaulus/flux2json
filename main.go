package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: flux2json [file]\n\n")
		fmt.Fprintf(os.Stderr, "Convert InfluxDB annotated CSV to JSON.\n")
		fmt.Fprintf(os.Stderr, "Reads from stdin when no file argument is given.\n")
	}
	flag.Parse()

	var input io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "flux2json: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	tables, err := Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "flux2json: parse error: %v\n", err)
		os.Exit(1)
	}

	output := struct {
		Tables []Table `json:"tables"`
	}{Tables: tables}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "flux2json: %v\n", err)
		os.Exit(1)
	}
}
