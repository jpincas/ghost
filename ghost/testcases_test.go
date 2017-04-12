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
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    "test",
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Single Where clause",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    "test",
				},
				WhereConfig{
					Key:      "id2",
					Operator: "=",
					Value:    "test2",
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test' AND id2 = 'test2') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Two Where clauses joined by AND",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    "test",
				},
				WhereConfig{
					Key:        "id2",
					Operator:   "=",
					Value:      "test2",
					JoinWithOr: true,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test' OR id2 = 'test2') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Two Where clauses joined by OR",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    "test",
				},
				WhereConfig{
					Key:      "id2",
					Operator: "=",
					Value:    "test2",
				},
				WhereConfig{
					Key:        "id3",
					Operator:   "=",
					Value:      "test3",
					JoinWithOr: true,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test' AND id2 = 'test2' OR id3 = 'test3') SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"3 Where clauses joined by AND then OR",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					AnyValue: []interface{}{1, 2, 3},
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = ANY(ARRAY[1, 2, 3])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Single multiple Value WHERE CLAUSE",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    "test",
				},
				WhereConfig{
					Key:        "name",
					AnyValue:   []interface{}{"jon", "jessi"},
					JoinWithOr: true,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = 'test' OR name = ANY(ARRAY['jon', 'jessi'])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Simple WHERE clause + multiple-Value any WHERE clause joined with OR",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    nil,
				},
				WhereConfig{
					Key:      "name",
					AnyValue: []interface{}{"jon", "jessi"},
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE name = ANY(ARRAY['jon', 'jessi'])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Nil WHERE clause + multiple-Value any WHERE clause",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    "",
				},
				WhereConfig{
					Key:      "name",
					AnyValue: []interface{}{"jon", "jessi"},
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE name = ANY(ARRAY['jon', 'jessi'])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Blank string WHERE clause + multiple-Value any WHERE clause",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "name",
					AnyValue: []interface{}{"jon", "jessi"},
				},
				WhereConfig{
					Key:      "id",
					Operator: "=",
					Value:    nil,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE name = ANY(ARRAY['jon', 'jessi'])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"Multiple-Value any WHERE clause + Nil WHERE clause",
	},
	{
		Query{
			Select: []string{"field1", "field2"},
			Schema: "public",
			Table:  "test_table",
			Where: []WhereConfig{
				WhereConfig{
					Key:      "id",
					AnyValue: []interface{}{1, 2, 3},
				},
				WhereConfig{
					Key:        "name",
					AnyValue:   []interface{}{"jon", "jessi"},
					JoinWithOr: true,
				},
			},
		},
		"WITH results AS (SELECT field1,field2 FROM public.test_table WHERE id = ANY(ARRAY[1, 2, 3]) OR name = ANY(ARRAY['jon', 'jessi'])) SELECT row_to_json(results) from results;",
		"[{'some':'object'}]",
		"2 x multiple-Value any WHERE clause joined with OR",
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
