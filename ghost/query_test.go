package ghost

import (
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {

	for _, c := range testCases {
		c.query.Build()
		if strings.ToLower(c.query.queryString) != strings.ToLower(c.expectedQueryString) {
			TestErrorFatal(t, c.description, c.query.queryString, c.expectedQueryString)
		}
	}

}
