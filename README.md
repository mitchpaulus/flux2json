# flux2json

Convert [InfluxDB annotated CSV](https://docs.influxdata.com/influxdb/v2/reference/syntax/annotated-csv/) query results into JSON.

## Installation

```bash
go install github.com/mitchpaulus/flux2json@latest
```

## Build from source

```bash
git clone https://github.com/mitchpaulus/flux2json.git
cd flux2json
go build -o flux2json .
```

## Usage

Read from a file:

```bash
flux2json results.csv
```

Read from stdin:

```bash
influx query 'from(bucket:"my-bucket") |> range(start: -1h)' --raw | flux2json
```

## Output format

```json
{
  "tables": [
    {
      "columns": [
        { "name": "_time",  "type": "datetime", "format": "rfc3339", "index": 0 },
        { "name": "host",   "type": "string",   "index": 1 },
        { "name": "_value", "type": "double",   "index": 2 }
      ],
      "data": [
        { "_time": "2026-02-26T12:00:00Z", "host": "a", "_value": 0.12 }
      ]
    }
  ]
}
```

Each annotated CSV table becomes one entry in the `tables` array. Column types are derived from the `#datatype` annotation and values are coerced to their native JSON types (`boolean`, `long`/`unsignedLong` → number, `double` → number, everything else → string).
