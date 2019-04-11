package plugins

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExtractJsonFromStruct_Vpc(t *testing.T) {
	fmt.Println(string(extractJsonFromStruct(VpcInput{})))
	fmt.Println(string(extractJsonFromStruct(VpcOutput{})))
}

func extractJsonFromStruct(input interface{}) string {
	output := ExtractJsonFromStruct(input)
	result, _ := json.MarshalIndent(output, "", "  ")
	return string(result)
}
