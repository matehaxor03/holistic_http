package http_extension

import (
	"fmt"
	"net/http"
	"strings"
	json "github.com/matehaxor03/holistic_json/json"
	common "github.com/matehaxor03/holistic_common/common"
)

func Nop() {
	
}

func WriteResponse(w http.ResponseWriter, result json.Map, write_response_errors []error) {
	
	if common.IsNil(write_response_errors) {
		var temp_errors []error
		write_response_errors = temp_errors
	}

	if common.IsNil(result) {
		temp_result := json.Map{}
		result = temp_result
	}
	
	result_errors, result_errors_errors := result.GetErrors("[errors]")
	
	if !common.IsNil(result_errors_errors) {
		write_response_errors = append(write_response_errors, result_errors_errors...)
	}

	if !common.IsNil(result_errors) {
		write_response_errors = append(write_response_errors, result_errors...)
	}

	if len(write_response_errors) > 0 {
		result.SetNil("data")
		result.SetErrors("[errors]", write_response_errors)
	} else {
		result.SetErrors("[errors]", nil)
	}

	var json_payload_builder strings.Builder
	result_as_string_errors := result.ToJSONString(&json_payload_builder)
	
	if result_as_string_errors != nil {
		write_response_errors = append(write_response_errors, result_as_string_errors...)
	}
	
	w.Header().Set("Content-Type", "application/json")

	if result_as_string_errors == nil {
		w.Write([]byte(json_payload_builder.String()))
	} else {
		w.Write([]byte(fmt.Sprintf("{\"unknown\":{\"[errors]\":\"%s\", \"data\":null}}", strings.ReplaceAll(fmt.Sprintf("%s", write_response_errors), "\"", "\\\""))))
	}
}
