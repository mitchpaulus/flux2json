---
title: Annotated CSV
description: InfluxDB and Flux return query results in annotated CSV format. You can also read annotated CSV directly from Flux with the csv.from() function, write data to InfluxDB using annotated CSV and the influx write command, or upload a CSV file in the UI.
url: https://docs.influxdata.com/influxdb/v2/reference/syntax/annotated-csv/
product: InfluxDB OSS v2
type: section
pages: 2
estimated_tokens: 12900
child_pages:
  - url: https://docs.influxdata.com/influxdb/v2/reference/syntax/annotated-csv/extended/
    title: Extended annotated CSV
---

# Annotated CSV

This page documents an earlier version of InfluxDB OSS. [InfluxDB 3 Core](/influxdb3/core/) is the latest stable version.

The InfluxDB `/api/v2/query` API returns query results in annotated CSV format. You can also write data to InfluxDB using annotated CSV and the `influx write` command, or [upload a CSV file](/influxdb/v2/write-data/no-code/load-data/) in the InfluxDB UI.

CSV tables must be encoded in UTF-8 and Unicode Normal Form C as defined in [UAX15](http://www.unicode.org/reports/tr15/). InfluxDB removes carriage returns before newline characters.

## CSV response format

InfluxDB annotated CSV supports encodings listed below.

### Tables

A table may have the following rows and columns.

#### Rows

-   **Annotation rows**: describe column properties.
-   **Header row**: defines column labels (one header row per table).
-   **Record row**: describes data in the table (one record per row).

##### Example

```csv
#group,false,false,true,true,false,false,true,true,true,true
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,double,string,string,string,string
#default,mean,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,host,region
,,0,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:52:00Z,15.43,mem,m,A,east
,,1,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:52:00Z,59.25,mem,m,B,east
,,2,2022-12-31T05:41:24Z,2023-01-31T05:41:24.001Z,2023-01-01T00:52:00Z,52.62,mem,m,C,east
```

#### Columns

In addition to the data columns, a table may include the following columns:

-   **Annotation column**: Displays the name of an annotation. Only used in annotation rows and is always the first column. Value can be empty or a supported [annotation](#annotations). The response format uses a comma (`,`) to separate an annotation name from values in the row. To account for this, rows in the table start with a leading comma; you’ll notice an empty column for the entire length of the table.
-   **Result column**: Contains the name of the result specified by the query.
-   **Table column**: Contains a unique ID for each table in a result.

### Multiple tables and results

If a file or data stream contains multiple tables or results, the following requirements must be met:

-   A table column indicates which table a row belongs to.
-   All rows in a table are contiguous.
-   An empty row delimits a new table boundary in the following cases:
    -   Between tables in the same result that do not share a common table schema.
    -   Between concatenated CSV files.
-   Each new table boundary starts with new annotation and header rows.

##### Example

```csv
#group,false,false,true,true,false,false,true,true,true,true
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
```

## Dialect options

Flux supports the following dialect options for `text/csv` format.

| Option | Description | Default |
| --- | --- | --- |
| header | If true, the header row is included. | true |
| delimiter | Character used to delimit columns. | , |
| quoteChar | Character used to quote values containing the delimiter. | " |
| annotations | List of annotations to encode (datatype, group, or default). | empty |
| commentPrefix | String prefix to identify a comment. Always added to annotations. | # |

## Annotations

Annotation rows describe column properties, and start with `#` (or commentPrefix value). The first column in an annotation row always contains the annotation name. Subsequent columns contain annotation values as shown in the table below.

| Annotation name | Values | Description |
| --- | --- | --- |
| datatype | a data type or line protocol element | Describes the type of data or which line protocol element the column represents. |
| group | boolean flag true or false | Indicates the column is part of the group key. |
| default | a value of the column’s data type | Value to use for rows with an empty value. |

Some tools might use or require a comma (`,`) to separate the annotation name from values in an annotation row.

To encode a table with its [group key](/influxdb/v2/reference/glossary/#group-key), the `datatype`, `group`, and `default` annotations must be included. If a table has no rows, the `default` annotation provides the group key values.

## Data types

| Datatype | Flux type | Description |
| --- | --- | --- |
| boolean | bool | “true” or “false” |
| unsignedLong | uint | unsigned 64-bit integer |
| long | int | signed 64-bit integer |
| double | float | IEEE-754 64-bit floating-point number |
| string | string | UTF-8 encoded string |
| base64Binary | bytes | base64 encoded sequence of bytes as defined in RFC 4648 |
| dateTime | time | instant in time, may be followed with a colon : and a description of the format (number, RFC3339, RFC3339Nano) |
| duration | duration | length of time represented as an unsigned 64-bit integer number of nanoseconds |

## Line protocol elements

The `datatype` annotation accepts [data types](#data-types) and **line protocol elements**. Line protocol elements identify how columns are converted into line protocol when using the [`influx write` command](/influxdb/v2/reference/cli/influx/write/) to write annotated CSV to InfluxDB.

| Line protocol element | Description |
| --- | --- |
| measurement | column value is the measurement |
| field (default) | column header is the field key, column value is the field value |
| tag | column header is the tag key, column value is the tag value |
| time | column value is the timestamp (alias for dateTime) |
| ignore orignored | column is ignored and not included in line protocol |

### Mixing data types and line protocol elements

Columns with [data types](#data-types) (other than `dateTime`) in the `#datatype` annotation are treated as **fields** when converted to line protocol. Columns without a specified data type default to `field` when converted to line protocol and **column values are left unmodified** in line protocol. *See an example [below](#example-of-mixing-data-types-line-protocol-elements) and [line protocol data types and format](/influxdb/v2/reference/syntax/line-protocol/#data-types-and-format).*

### Time columns

A column with `time` or `dateTime` `#datatype` annotations are used as the timestamp when converted to line protocol. If there are multiple `time` or `dateTime` columns, the last column (on the right) is used as the timestamp in line protocol. Other time columns are ignored and the `influx write` command outputs a warning.

Time column values should be **Unix nanosecond timestamps**, **RFC3339**, or **RFC3339Nano**.

##### Example line protocol elements in datatype annotation

```csv
#group false,false,false,false,false,false,false
#datatype measurement,tag,tag,field,field,ignored,time
#default ,,,,,,
m,cpu,host,time_steal,usage_user,nothing,time
cpu,cpu1,host1,0,2.7,a,1482669077000000000
cpu,cpu1,host2,0,2.2,b,1482669087000000000
```

Resulting line protocol:

```text
cpu,cpu=cpu1,host=host1 time_steal=0,usage_user=2.7 1482669077000000000
cpu,cpu=cpu1,host=host2 time_steal=0,usage_user=2.2 1482669087000000000
```

##### Example of mixing data types and line protocol elements

```csv
#group,false,false,false,false,false,false,false,false,false
#datatype,measurement,tag,string,double,boolean,long,unsignedLong,duration,dateTime
#default,test,annotatedDatatypes,,,,,,
,m,name,s,d,b,l,ul,dur,time
,,,str1,1.0,true,1,1,1ms,1
,,,str2,2.0,false,2,2,2us,2020-01-11T10:10:10Z
```

Resulting line protocol:

```text
test,name=annotatedDatatypes s="str1",d=1,b=true,l=1i,ul=1u,dur=1000000i 1
test,name=annotatedDatatypes s="str2",d=2,b=false,l=2i,ul=2u,dur=2000i 1578737410000000000
```

## Errors

If an error occurs during execution, a table returns with:

-   An error column that contains an error message.
-   A reference column with a unique reference code to identify more information about the error.
-   A second row with error properties.

If an error occurs:

-   Before InfluxDB responds with data, the HTTP status code indicates an error. The output table contains error details.
-   After InfluxDB sends partial results to the client, the error is encoded as the next table and remaining results are discarded. In this case, the HTTP status code remains 200 OK.

##### Example

Encoding for an error with the datatype annotation:

```csv
#datatype,string,long
,error,reference
,Failed to parse query,897
```

#### Related

-   [csv.from() function](/flux/v0/stdlib/csv/from/)
-   [Extended annotated CSV](/influxdb/v2/reference/syntax/annotated-csv/extended/)

[csv](/influxdb/v2/tags/csv/) [syntax](/influxdb/v2/tags/syntax/)


---

## Extended annotated CSV

This page documents an earlier version of InfluxDB OSS. [InfluxDB 3 Core](/influxdb3/core/) is the latest stable version.

**Extended annotated CSV** provides additional annotations and options that specify how CSV data should be converted to [line protocol](/influxdb/v2/reference/syntax/line-protocol/) and written to InfluxDB. InfluxDB uses the [`csv2lp` library](https://github.com/influxdata/influxdb/tree/master/pkg/csv2lp) to convert CSV into line protocol. Extended annotated CSV supports all [Annotated CSV](/influxdb/v2/reference/syntax/annotated-csv/) annotations.

The Flux [`csv.from` function](/flux/v0/stdlib/csv/from/) only supports [annotated CSV](/influxdb/v2/reference/syntax/annotated-csv/), not extended annotated CSV.

To write data to InfluxDB, line protocol must include the following:

-   [measurement](/influxdb/v2/reference/syntax/line-protocol/#measurement)
-   [field set](/influxdb/v2/reference/syntax/line-protocol/#field-set)
-   [timestamp](/influxdb/v2/reference/syntax/line-protocol/#timestamp) *(Optional but recommended)*
-   [tag set](/influxdb/v2/reference/syntax/line-protocol/#tag-set) *(Optional)*

Extended CSV annotations identify the element of line protocol a column represents.

## CSV Annotations

Extended annotated CSV extends and adds the following annotations:

-   [datatype](#datatype)
-   [constant](#constant)
-   [timezone](#timezone)
-   [concat](#concat)

### datatype

Use the `#datatype` annotation to specify the [line protocol element](/influxdb/v2/reference/syntax/line-protocol/#elements-of-line-protocol) a column represents. To explicitly define a column as a **field** of a specific data type, use the field type in the annotation (for example: `string`, `double`, `long`, etc.).

| Data type | Resulting line protocol |
| --- | --- |
| measurement | Column is the measurement |
| tag | Column is a tag |
| dateTime | Column is the timestamp |
| field | Column is a field |
| ignored | Column is ignored |
| string | Column is a string field |
| double | Column is a float field |
| long | Column is an integer field |
| unsignedLong | Column is an unsigned integer field |
| boolean | Column is a boolean field |

#### measurement

Indicates the column is the **measurement**.

#### tag

Indicates the column is a **tag**. The **column label** is the **tag key**. The **column value** is the **tag value**.

#### dateTime

Indicates the column is the **timestamp**. `time` is an alias for `dateTime`. If the [timestamp format](#supported-timestamp-formats) includes a time zone, the parsed timestamp includes the time zone offset. By default, all timestamps are UTC. You can also use the [`#timezone` annotation](#timezone) to adjust timestamps to a specific time zone.

There can only be **one** `dateTime` column.

The `influx write` command converts timestamps to [Unix timestamps](/influxdb/v2/reference/glossary/#unix-timestamp). Append the timestamp format to the `dateTime` datatype with (`:`).

```csv
#datatype dateTime:RFC3339
#datatype dateTime:RFC3339Nano
#datatype dateTime:number
#datatype dateTime:2006-01-02
```

##### Supported timestamp formats

| Timestamp format | Description | Example |
| --- | --- | --- |
| RFC3339 | RFC3339 timestamp | 2020-01-01T00:00:00Z |
| RFC3339Nano | RFC3339 timestamp | 2020-01-01T00:00:00.000000000Z |
| number | Unix timestamp | 1577836800000000000 |

If using the `number` timestamp format and timestamps are **not in nanoseconds**, use the [`influx write --precision` flag](/influxdb/v2/reference/cli/influx/write/#flags) to specify the [timestamp precision](/influxdb/v2/reference/glossary/#precision).

##### Custom timestamp formats

To specify a custom timestamp format, use timestamp formats as described in the [Go time package](https://golang.org/pkg/time). For example: `2020-01-02`.

#### field

Indicates the column is a **field**. The **column label** is the **field key**. The **column value** is the **field value**.

With the `field` datatype, field values are copies **as-is** to line protocol. For information about line protocol values and how they are written to InfluxDB, see [Line protocol data types and formats](/influxdb/v2/reference/syntax/line-protocol/#data-types-and-format). We generally recommend specifying the [field type](#field-types) in annotations.

#### ignored

The column is ignored and not written to InfluxDB.

#### Field types

The column is a **field** of a specified type. The **column label** is the **field key**. The **column value** is the **field value**.

-   [string](#string)
-   [double](#double)
-   [long](#long)
-   [unsignedLong](#unsignedlong)
-   [boolean](#boolean)

##### string

Column is a **[string](/influxdb/v2/reference/glossary/#string) field**.

##### double

Column is a **[float](/influxdb/v2/reference/glossary/#float) field**. By default, InfluxDB expects float values that use a period (`.`) to separate the fraction from the whole number. If column values include or use other separators, such as commas (`,`) to visually separate large numbers into groups, specify the following **float separators**:

-   **fraction separator**: Separates the fraction from the whole number.
-   **ignored separator**: Visually separates the whole number into groups but ignores the separator when parsing the float value.

Use the following syntax to specify **float separators**:

```sh
# Syntax
<fraction-separator><ignored-separator>

# Example
.,

# With the float separators above
# 1,200,000.15 => 1200000.15
```

Append **float separators** to the `double` datatype annotation with a colon (`:`). For example:

```
#datatype "double:.,"
```

If your **float separators** include a comma (`,`), wrap the column annotation in double quotes (`""`) to prevent the comma from being parsed as a column separator or delimiter. You can also [define a custom column separator](#define-custom-column-separator).

##### long

Column is an **[integer](/influxdb/v2/reference/glossary/#integer) field**. If column values contain separators such as periods (`.`) or commas (`,`), specify the following **integer separators**:

-   **fraction separator**: Separates the fraction from the whole number. ***Integer values are truncated at the fraction separator when converted to line protocol.***
-   **ignored separator**: Visually separates the whole number into groups but ignores the separator when parsing the integer value.

Use the following syntax to specify **integer separators**:

```sh
# Syntax
<fraction-separator><ignored-separator>

# Example
.,

# With the integer separators above
# 1,200,000.00 => 1200000i
```

Append **integer separators** to the `long` datatype annotation with a colon (`:`). For example:

```
#datatype "long:.,"
```

If your **integer separators** include a comma (`,`), wrap the column annotation in double quotes (`""`) to prevent the comma from being parsed as a column separator or delimiter. You can also [define a custom column separator](#define-custom-column-separator).

##### unsignedLong

Column is an **[unsigned integer (uinteger)](/influxdb/v2/reference/glossary/#unsigned-integer) field**. If column values contain separators such as periods (`.`) or commas (`,`), specify the following **uinteger separators**:

-   **fraction separator**: Separates the fraction from the whole number. ***Uinteger values are truncated at the fraction separator when converted to line protocol.***
-   **ignored separator**: Visually separates the whole number into groups but ignores the separator when parsing the uinteger value.

Use the following syntax to specify **uinteger separators**:

```sh
# Syntax
<fraction-separator><ignored-separator>

# Example
.,

# With the uinteger separators above
# 1,200,000.00 => 1200000u
```

Append **uinteger separators** to the `long` datatype annotation with a colon (`:`). For example:

```
#datatype "unsignedLong:.,"
```

If your **uinteger separators** include a comma (`,`), wrap the column annotation in double quotes (`""`) to prevent the comma from being parsed as a column separator or delimiter. You can also [define a custom column separator](#define-custom-column-separator).

##### boolean

Column is a **[boolean](/influxdb/v2/reference/glossary/#boolean) field**. If column values are not [supported boolean values](/influxdb/v2/reference/syntax/line-protocol/#boolean), specify the **boolean format** with the following syntax:

```sh
# Syntax
<true-values>:<false-values>

# Example
y,Y,1:n,N,0

# With the boolean format above
# y => true, Y => true, 1 => true
# n => false, N => false, 0 => false
```

Append the **boolean format** to the `boolean` datatype annotation with a colon (`:`). For example:

```
#datatype "boolean:y,Y:n,N"
```

If your **boolean format** contains commas (`,`), wrap the column annotation in double quotes (`""`) to prevent the comma from being parsed as a column separator or delimiter. You can also [define a custom column separator](#define-custom-column-separator).

### constant

Use the `#constant` annotation to define a constant column label and value for each row. The `#constant` annotation provides a way to supply [line protocol elements](/influxdb/v2/reference/syntax/line-protocol/#elements-of-line-protocol) that don’t exist in the CSV data.

Use the following syntax to define constants:

```
#constant <datatype>,<column-label>,<column-value>
```

To provide multiple constants, include each `#constant` annotations on a separate line.

```
#constant measurement,m
#constant tag,dataSource,csv
```

For constants with `measurement` and `dateTime` datatypes, the second value in the constant definition is the **column-value**.

### timezone

Use the `#timezone` annotation to update timestamps to a specific timezone. By default, timestamps are parsed as UTC. Use the `±HHmm` format to specify the timezone offset relative to UTC.

### strict mode

Use the `:strict` keyword to indicate a loss of precision when parsing `long` or `unsignedLong` data types. Turn on strict mode by using a column data type that ends with `strict`, such as `long:strict`. When parsing `long` or `unsignedLong` value from a string value with fraction digits, the whole CSV row fails when in a strict mode. A warning is printed when not in a strict mode, saying `line x: column y: '1.2' truncated to '1' to fit into long data type`. For more information on strict parsing, see the [package documentation](https://github.com/influxdata/influxdb/tree/master/pkg/csv2lp).

##### Timezone examples

| Timezone | Offset |
| --- | --- |
| US Mountain Daylight Time | -0600 |
| Central European Summer Time | +0200 |
| Australia Eastern Standard Time | +1000 |
| Apia Daylight Time | +1400 |

##### Timezone annotation example

```
#timezone -0600
```

### concat

The `#concat` annotation adds a new column that is concatenated from existing columns according to bash-like string interpolation literal with variables referencing existing column labels.

For example:

```
#concat,string,fullName,${firstName} ${lastName}
```

This is especially useful when constructing a timestamp from multiple columns. For example, the following annotation will combine the given CSV columns into a timestamp:

```
#concat,dateTime:2006-01-02,${Year}-${Month}-${Day}

Year,Month,Day,Hour,Minute,Second,Tag,Value
2020,05,22,00,00,00,test,0
2020,05,22,00,05,00,test,1
2020,05,22,00,10,00,test,2
```

## Define custom column separator

If columns are delimited using a character other than a comma, use the `sep` keyword to define a custom separator **in the first line of your CSV file**.

```
sep=;
```

## Annotation shorthand

Extended annotated CSV supports **annotation shorthand**. Include the column label, datatype, and *(optional)* default value in each column header row using the following syntax:

```
<column-label>|<column-datatype>|<column-default-value>
```

##### Example annotation shorthand

```
m|measurement,location|tag|Hong Kong,temp|double,pm|long|0,time|dateTime:RFC3339
weather,San Francisco,51.9,38,2020-01-01T00:00:00Z
weather,New York,18.2,,2020-01-01T00:00:00Z
weather,,53.6,171,2020-01-01T00:00:00Z
```

##### The shorthand explained

-   The `m` column represents the **measurement** and has no default value.
-   The `location` column is a **tag** with the default value, `Hong Kong`.
-   The `temp` column is a **field** with **float** (`double`) values and no default value.
-   The `pm` column is a **field** with **integer** (`long`) values and a default of `0`.
-   The `time` column represents the **timestamp**, uses the **RFC3339** timestamp format, and has no default value.

##### Resulting line protocol

```
weather,location=San\ Francisco temp=51.9,pm=38i 1577836800000000000
weather,location=New\ York temp=18.2,pm=0i 1577836800000000000
weather,location=Hong\ Kong temp=53.6,pm=171i 1577836800000000000
```

#### Related

-   [Write CSV data to InfluxDB](/influxdb/v2/write-data/developer-tools/csv/)
-   [influx write](/influxdb/v2/reference/cli/influx/write/)
-   [Line protocol](/influxdb/v2/reference/syntax/line-protocol/)
-   [Annotated CSV](/influxdb/v2/reference/syntax/annotated-csv/)

[csv](/influxdb/v2/tags/csv/) [syntax](/influxdb/v2/tags/syntax/) [write](/influxdb/v2/tags/write/)
