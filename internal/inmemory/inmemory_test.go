package inmemory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInmemory(t *testing.T) {
	c := New()

	tests := []struct {
		name     string
		value    string
		key      string
		wantGet  string
		wantSet  bool
		foundKey string
	}{
		{
			name:     "shorten",
			value:    "value1",
			key:      "key1",
			foundKey: "key1",
			wantGet:  "value1",
			wantSet:  true,
		},
		{
			name:     "duplicate record",
			value:    "value1",
			key:      "key1",
			foundKey: "key1",
			wantGet:  "value1",
			wantSet:  false,
		},
		{
			name:     "not value",
			value:    "value2",
			key:      "key2",
			foundKey: "key3",
			wantGet:  "",
			wantSet:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			update, _ := c.CreateOrUpdate(test.key, test.value)
			assert.Equal(t, test.wantSet, update)
			value, _ := c.GetValue(test.foundKey)
			assert.Equal(t, test.wantGet, value)
		})
	}
}
