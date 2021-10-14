package bitbucket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const outputStr = `
some
other
text
--- OUTPUT JSON START ---
{
	"key_output": "value",
	"map_output": {
		"map_key": "map_value"
	}
}
--- OUTPUT JSON STOP ---

`

func TestOutputBlock(t *testing.T) {
	d, err := extractOutputs(outputStr)
	assert.Nil(t, err, "err should be nil")
	assert.Containsf(t, d, "key_output", "expected key_output")
	assert.Equal(t, d["key_output"], "value", "expected value")
	assert.Containsf(t, d, "map_output", "expected map_output")
	assert.Equal(t, d["map_output"], map[string]interface{}{"map_key": "map_value"}, "expected map_key with map_value")
}
