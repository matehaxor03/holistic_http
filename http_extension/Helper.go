package http_extension

import (
	"fmt"
	"net/http"
	"strings"
	json "github.com/matehaxor03/holistic_json/json"
)

func Nop() {
	
}

func WriteResponse(w http.ResponseWriter, result json.Map, write_response_errors []error) {
	keys := result.Keys()
	
	if len(keys) != 1 {
		write_response_errors = append(write_response_errors, fmt.Errorf(fmt.Sprintf("number of root keys is incorrect %s",keys)))
	}
	
	
	inner_map_found := false
	if len(keys) == 1 {
		inner_map, inner_map_errors := result.GetMap(keys[0])
		if inner_map_errors != nil {
			write_response_errors = append(write_response_errors, inner_map_errors...)
		} 
		
		if inner_map == nil {
			write_response_errors = append(write_response_errors, fmt.Errorf("inner map is nil"))
			inner_map_found = false
		} else {
			inner_map_found = true
		}
	}


	if len(write_response_errors) > 0 {
		if inner_map_found {
			(result[keys[0]].(json.Map))["data"] = nil
			(result[keys[0]].(json.Map))["[errors]"] = write_response_errors
		} else {
			result["unknown"] = json.Map{"data":nil, "[errors]":write_response_errors}
		}
	} else {
		if inner_map_found {
			(result[keys[0]].(json.Map))["[errors]"] = write_response_errors
		} else {
			result["unknown"] = json.Map{"data":nil, "[errors]":write_response_errors}
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
		w.Write([]byte(fmt.Sprintf("{\"unknown\":{\"[errors]\":\"%s\", \"data\":null}}", strings.ReplaceAll(fmt.Sprintf("%s", result_as_string_errors), "\"", "\\\""))))
	}
}
