package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	t.Parallel()

	tests := []any{"string", 123, 12.3, true}

	for _, test := range tests {
		assert.Equal(t, test, *Ptr(test))
	}
}
