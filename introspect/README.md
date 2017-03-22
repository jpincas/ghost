# GO PostgreSQL Table/View to Schema Struct and/or JSON

Given a PostgreSQL table or view, this package will create a Go struct describing the schema structure of the table, and optionally output the schema to JSON.

The schema is basically a list (Go map, JSON object) keyed by DB column names, with attributes describing the properties of each column, namely its data type, defaults and restrictions.


## Why would you need this?

Automatically generating schemas for your tables/views might be useful if:

- You are building an API in Go using PostgreSQL and there is a 1:1 relationship between the table (or more likely view) and the JSON your API outputs, this package will basically give you a JSON schema describing your API.
- You'd like to autogenerate lists or forms in templates or client apps based on the structure of your database tables.

I wrote this package to provide the JSON schema functionality for [EcoSystem](https://github.com/ecosystemsoftware/ecosystem), which uses JSON output directly from PostgreSQL and pipes it through to a JSON API endpoint with zero manipulation - thus the JSON schema of a table perfectly describes the API endpoint for this table.

## Usage

Simply call either function with the database connection pointer `core.DB` , the name of your database schema `dbSchema` ('public' if you're not using schemas), the name of the table or view `dbTable`, and the name of the databse role to set e.g.,

```
s, err := GetSchema(db, dbSchema, dbTable, "admin")
json, err := GetSchemaJSON(db, dbSchema, dbTable, "admin")
```

Use it in your handlers to output JSON in your API, or for further processing.




### Data Types Supported

| PG Data Type               | Schema Type |
| -------------------------- | ----------- |
| numeric, bigint, integer   | number      |
| array                      | array       |
| timestamp without timezone | date        |
| boolean                    | boolean     |
| DEFAULT                    | string      |



### PostgreSQL Modifiers Supported

| PG Modifier              | Schema Attribute                       |
| ------------------------ | -------------------------------------- |
| NOT NULL                 | required: true                         |
| DEFAULT x                | default: x                             |
| VARCHAR(x)               | maxlength: x                           |
| REFERENCES table(column) | reftable: "table"; refcolumn: "column" |