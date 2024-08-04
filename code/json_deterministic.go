package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/gowebpki/jcs"
)

const prefix = "_id:"

// Function to compute SHA-256 hash of a string
func computeHash(value []byte) string {
	hash := sha256.New()
	hash.Write(value)
	return prefix + hex.EncodeToString(hash.Sum(nil))
}

// Function to walk through the JSON and detect tables to transform them
// Issues
// hash is hash of the content -> could be used to brute force the information
// Use the hash only to sort the elements, then replace them by indexes 0, 1, 2, ...

//

func walkAndTransformJSON(input interface{}) interface{} {
	switch value := input.(type) {
	// List
	case []interface{}:
		new_map := make(map[string]interface{})
		index_list := make(map[string]int)
		for _, item := range value {
			transformed := walkAndTransformJSON(item)
			transformedByte, err := json.Marshal(transformed)
			if err != nil {
				println("ERROR", err.Error())
			}

			index := computeHash(transformedByte)
			// Check if index exists
			_, ok := index_list[index]
			// If the key exists
			if ok {
				_index := index
				index_list[_index] += 1
				index = fmt.Sprintf("%s_%d", _index, index_list[_index])
			} else {
				index_list[index] = 0
			}

			new_map[index] = transformed
		}
		return new_map

	// Object
	case map[string]interface{}:
		valueByte, err := json.Marshal(value)
		if err != nil {
			println("ERROR: ", err.Error())
		}
		valueByteJCS, err := jcs.Transform(valueByte)
		if err != nil {
			println("ERROR: ", err.Error())
		}
		err = json.Unmarshal(valueByteJCS, &value)
		if err != nil {
			println("ERROR: ", err.Error())
		}
		// Recursively process each value in the map
		for k, v := range value {
			value[k] = walkAndTransformJSON(v)
		}
	}
	return input
}

func walkAndReplace(input interface{}) interface{} {
	switch value := input.(type) {
	// Object
	case map[string]interface{}:
		result := make(map[string]interface{})
		keys := []string{}
		for k := range value {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// Recursively process each value in the map
		index := 0
		for _, k := range keys {
			v := value[k]
			if strings.HasPrefix(k, prefix) {
				id := fmt.Sprintf("%d", index)
				result[id] = walkAndReplace(v)
				index += 1
			} else {
				result[k] = walkAndReplace(v)
			}
		}
		return result
	}
	return input
}

// Function to walk through the JSON and detect tables to transform them
func walkAndTransformToJSONPointer(input interface{}, path string, output map[string]interface{}) interface{} {
	switch value := input.(type) {
	// Object
	case map[string]interface{}:
		// Recursively process each value in the map
		for k, v := range value {
			_path := fmt.Sprintf("%s/%s", path, k)
			walkAndTransformToJSONPointer(v, _path, output)
		}
	default:
		// fmt.Println(path, input)
		output[path] = input
	}
	return ""
}

func jsonPointerToArray(input map[string]interface{}) []string {
	var result []string
	for k, value := range input {
		var stringValue string
		switch v := value.(type) {
		case string:
			stringValue = v
		case float64:
			stringValue = fmt.Sprintf("%f", v)
		case bool:
			stringValue = fmt.Sprintf("%t", v)
		case nil:
			stringValue = "null"
		default:
			stringValue = fmt.Sprintf("%v", v)
		}

		item := fmt.Sprintf("%s:%s", k, stringValue)
		result = append(result, item)
	}

	sort.Strings(result)
	return result
}

func arrayToScalar(input []string) []curves.Scalar {

	curve := curves.BLS12381(&curves.PointBls12381G2{})

	var msg []curves.Scalar
	for _, element := range input {
		msg = append(msg, curve.Scalar.Hash([]byte(element)))
	}

	return msg
}

// func main() {
// 	jsonInput := []byte(`{
//         "users": [
//             {"name": "Charlie", "age": 35, "active": null},
//             {"name": "Bob", "age": 25, "active": [1, 1, 1]},
//             {"name": "Alice", "age": 30, "active": [true, true, true]}
//         ],
//         "data": {
//             "info": [
//                 {"id": 1, "value": "x"},
//                 {"id": 1, "value": "x"}
//             ]
//         },
//         "emptyArray": [],
//         "emptyObject": {},
//         "stringValue": "example",
//         "numberValue": 123.456,
//         "booleanValue": true,
//         "nullValue": null
//     }`)
//
// 	var jsonData interface{}
// 	err := json.Unmarshal(jsonInput, &jsonData)
// 	if err != nil {
// 		log.Fatalf("Error unmarshalling JSON: %v", err)
// 	}
//
// 	transformedData := walkAndTransformJSON(jsonData)
//
// 	transformedJSON, err := json.MarshalIndent(transformedData, "", "  ")
// 	if err != nil {
// 		log.Fatalf("Error marshalling transformed JSON: %v", err)
// 	}
//
// 	fmt.Println(string(transformedJSON))
//
// 	jsonPointer := make(map[string]interface{})
//
// 	_ = walkAndTransformToJSONPointer(transformedData, "", jsonPointer)
// 	output, err := json.MarshalIndent(jsonPointer, "", "  ")
// 	if err != nil {
// 		log.Fatalf("Error marshalling transformed JSON: %v", err)
// 	}
//
// 	fmt.Println(string(output))
// }
