package bitbucket

import (
	"encoding/json"
	"fmt"
	"strings"
)

const jsonStart = "--- OUTPUT JSON START ---"
const jsonEnd = "--- OUTPUT JSON STOP ---"

// extractOutputs of a bitbucket log, as output are not supported by bitbucket itself
func extractOutputs(log string) (result map[string]interface{}, err error) {
	start := strings.Index(log, jsonStart)
	stop := strings.Index(log, jsonEnd)
	if start != -1 && stop != -1{
		dataStr := log[start+len(jsonStart) : stop]
		err = json.Unmarshal([]byte(dataStr), &result)
	} else {
		fmt.Println("[DEBUG]: no outputs found")
	}
	return
}
