package ghost

import "testing"

func TestBuild(t *testing.T) {

	for _, c := range testCases {
		c.query.Build()
		if c.query.queryString != c.expectedQueryString {
			TestErrorFatal(t, c.description, c.query.queryString, c.expectedQueryString)
		}
	}

}
