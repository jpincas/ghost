package ghost

var testCases = []struct {
	query               Query
	expectedQueryString string
	mockResult          string
	description         string
}{
	{
		Query{
			BaseSQL: "SELECT * FROM %s.%s",
			SQLArgs: []interface{}{"public", "test_table"},
		},
		"WITH results AS (SELECT * FROM public.test_table) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Base SQL + Args",
	},
	{
		Query{
			Select: []string{"*"},
			Schema: "public",
			Table:  "test_table",
		},
		"WITH results AS (SELECT * FROM public.test_table) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Select specified with schema and table",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Multiple select fields",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []whereConfig{
				whereConfig{
					key:      "id",
					operator: "=",
					value:    "test",
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Single where clause",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []whereConfig{
				whereConfig{
					key:      "id",
					operator: "=",
					value:    "test",
				},
				whereConfig{
					key:      "id2",
					operator: "=",
					value:    "test2",
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test' AND id2 = 'test2') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Two where clauses joined by AND",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []whereConfig{
				whereConfig{
					key:      "id",
					operator: "=",
					value:    "test",
				},
				whereConfig{
					key:        "id2",
					operator:   "=",
					value:      "test2",
					joinWithOr: true,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test' OR id2 = 'test2') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Two where clauses joined by OR",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []whereConfig{
				whereConfig{
					key:      "id",
					anyValue: []interface{}{1, 2, 3},
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = ANY(ARRAY[1, 2, 3])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Single multiple value WHERE CLAUSE",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []whereConfig{
				whereConfig{
					key:      "id",
					anyValue: []interface{}{1, 2, 3},
				},
				whereConfig{
					key:        "name",
					anyValue:   []interface{}{"jon", "jessi"},
					joinWithOr: true,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = ANY(ARRAY[1, 2, 3]) OR name = ANY(ARRAY['jon', 'jessi'])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Simple WHERE clause + multiple-value any WHERE clause joined with OR",
	},
	{
		Query{
			Select: []string{"*"},
			Schema: "public",
			Table:  "test_table",
			IsList: true,
		},
		"WITH results AS (SELECT * FROM public.test_table) SELECT array_to_json(array_agg(row_to_json(results))) from results;",
		"[{'some':'object'}]",
		"Select specified with schema and table, return a list",
	},
	{
		Query{
			Select: []string{"*"},
			Schema: "public",
			Table:  "test_table",
			IsList: true,
			Role:   "admin",
		},
		"SET LOCAL ROLE admin; WITH results AS (SELECT * FROM public.test_table) SELECT array_to_json(array_agg(row_to_json(results))) from results;",
		"[{'some':'object'}]",
		"Select specified with schema and table, return a list, add role",
	},
	{
		Query{
			Select: []string{"*"},
			Schema: "public",
			Table:  "test_table",
			IsList: true,
			Role:   "admin",
			UserID: "123456",
		},
		"SET my.user_id = '123456'; SET LOCAL ROLE admin; WITH results AS (SELECT * FROM public.test_table) SELECT array_to_json(array_agg(row_to_json(results))) from results;",
		"[{'some':'object'}]",
		"Select specified with schema and table, return a list, add role and user id",
	},
}

// {
// 		Query{
// 			Select: []string{"*"},
// 			Schema: "public",
// 			Table:  "test_table",
// 			UserID: "123456",
// 			IsList: true,
// 		},
// 		"SET my.user_id = '123456'; WITH results AS (SELECT * FROM public.test_table) SELECT array_to_json(array_agg(row_to_json(results))) from results;",
// 		"[{'some':'object'}]",
// 		"Select specified with schema and table, return a list, user id set",
// 	},
