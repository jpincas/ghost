// Copyright 2017 Jonathan Pincas

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"context"

	"net/http"

	"github.com/pressly/chi"
)

func AddSchemaAndTableToContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "schema", HyphensToUnderscores(chi.URLParam(r, "schema")))
		ctx = context.WithValue(ctx, "table", HyphensToUnderscores(chi.URLParam(r, "table")))
		ctx = context.WithValue(ctx, "queries", r.URL.Query())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AddRecordToContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "record", chi.URLParam(r, "record"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
