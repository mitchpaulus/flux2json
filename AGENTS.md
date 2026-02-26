This is a repository for a simple CLI utility with one purpose: Convert the structured CSV output from InfluxDb into JSON.

Default schema (There may be others)

```
{
  "tables": [
     {
       "columns": [
          { "name": "_time", "type": "datetime", "format": "rfc3339", "index": 0 },
          { "name": "host", "type": "string", "index": 1 },
          { "name": "cpu", "type": "float", "index": 2 },
          { "name": "ok", "type": "bool", "index": 3 }
       ]
       "data": [
            {
                "_time": "2026-02-26T12:00:00Z",
                "host": "a",
                "cpu": 0.12,
                "ok": true
            },
              ...
        ]
     }
  ]
}
```

Please always read `annotated-csv.md` as this provide all the documentation about the annotated CSV format.
