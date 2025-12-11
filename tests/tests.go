package tests

import (
	"strings"
	"testing"

	. "github.com/aprikotdev/speckles/elements"
	"github.com/stretchr/testify/assert"
)

type result struct {
	Expected string
	Actual   ElementRenderer
}

func run(t *testing.T, results []result) {
	for _, result := range results {
		var sb strings.Builder
		e := result.Expected

		err := result.Actual.Render(&sb)
		assert.NoError(t, err)

		a := sb.String()
		assert.Equal(t, e, a)
	}
}
