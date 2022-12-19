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
	keys := result.Keys()

	if common.IsNil(write_response_errors) {
		var temp_errors []error
		write_response_errors = temp_errors
	}
	
	if len(keys) != 1 {
		write_response_errors = append(write_response_errors, fmt.Errorf(fmt.Sprintf("number of root keys is incorrect %s",keys)))
	}
	
	var inner_map_value *json.Map
	inner_map_found := false
	if len(keys) == 1 {
		inner_map, inner_map_errors := result.GetMap(keys[0])
		if inner_map_errors != nil {
			write_response_errors = append(write_response_errors, inner_map_errors...)
		} else if common.IsNil(inner_map) {
			write_response_errors = append(write_response_errors, fmt.Errorf("inner map is nil"))
			inner_map_found = false
		} else {
			inner_map_found = true
			inner_map_value = inner_map

			inner_map_errors, inner_map_errors_errors := inner_map_value.GetErrors("[errors]")
			if inner_map_errors_errors != nil {
				write_response_errors = append(write_response_errors, inner_map_errors_errors...)
			} else if inner_map_errors != nil {
				write_response_errors = append(write_response_errors, inner_map_errors...)
			}
		}
	}


	if len(write_response_errors) > 0 {
		if inner_map_found {
			inner_map_value.SetNil("data")
			inner_map_value.SetErrors("[errors]", &write_response_errors)
		} else {
			result["unknown"] = json.Map{"data":nil, "[errors]":write_response_errors}
		}
	} else {
		if inner_map_found {
			inner_map_value.SetErrors("[errors]", nil)
		} else {
			result["unknown"] = json.Map{"data":nil, "[errors]":nil}
		}
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
