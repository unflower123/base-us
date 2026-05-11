package mapx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMapToJson(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	jsonStr, err := MapToJson(data)
	if err != nil {
		t.Fatalf("MapToJson failed: %v", err)
	}
	expected := `{"a":1,"b":2}`
	if jsonStr != expected {
		t.Errorf("MapToJson returned %s, expected %s", jsonStr, expected)
	}
}

func TestJsonToMap(t *testing.T) {
	jsonStr := `{"a":1,"b":2}`
	data, err := JsonToMap[string, int](jsonStr)
	if err != nil {
		t.Fatalf("JsonToMap failed: %v", err)
	}
	expected := map[string]int{
		"a": 1,
		"b": 2,
	}
	for k, v := range expected {
		if data[k] != v {
			t.Errorf("JsonToMap returned %v, expected %v", data, expected)
		}
	}
}

func TestStructToMap(t *testing.T) {
	type TestStruct struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	obj := TestStruct{A: 1, B: "test"}
	data, err := StructToMap(obj)
	if err != nil {
		t.Fatalf("StructToMap failed: %v", err)
	}
	expected := map[string]interface{}{
		"a": float64(1), // JSON unmarshal converts numbers to float64
		"b": "test",
	}
	for k, v := range expected {
		if data[k] != v {
			t.Errorf("StructToMap returned %v, expected %v", data, expected)
		}
	}
}

func TestMapToStruct(t *testing.T) {
	type TestStruct struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	data := map[string]interface{}{
		"a": 1,
		"b": "test",
	}
	var obj TestStruct
	err := MapToStruct(data, &obj)
	if err != nil {
		t.Fatalf("MapToStruct failed: %v", err)
	}
	expected := TestStruct{A: 1, B: "test"}
	if obj != expected {
		t.Errorf("MapToStruct returned %v, expected %v", obj, expected)
	}
}

func TestMergeMaps(t *testing.T) {
	map1 := map[string]int{
		"a": 1,
		"b": 2,
	}
	map2 := map[string]int{
		"b": 3,
		"c": 4,
	}
	merged := MergeMaps(map1, map2, true)
	expected := map[string]int{
		"a": 1,
		"b": 3,
		"c": 4,
	}
	for k, v := range expected {
		if merged[k] != v {
			t.Errorf("MergeMaps returned %v, expected %v", merged, expected)
		}
	}
}

func TestFilterMap(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	filtered := FilterMap(data, func(key string, value int) bool {
		return value > 1
	})
	expected := map[string]int{
		"b": 2,
		"c": 3,
	}
	for k, v := range expected {
		if filtered[k] != v {
			t.Errorf("FilterMap returned %v, expected %v", filtered, expected)
		}
	}
}

func TestMapValues(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	mapped := MapValues(data, func(key string, value int) string {
		return fmt.Sprintf("%s:%d", key, value)
	})
	expected := map[string]string{
		"a": "a:1",
		"b": "b:2",
	}
	for k, v := range expected {
		if mapped[k] != v {
			t.Errorf("MapValues returned %v, expected %v", mapped, expected)
		}
	}
}

func TestMapKeysToLower(t *testing.T) {
	data := map[string]interface{}{
		"A": 1,
		"B": map[string]interface{}{
			"C": 2,
		},
	}
	lower := MapKeysToLower(data)
	expected := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": 2,
		},
	}
	if !reflect.DeepEqual(lower, expected) {
		t.Errorf("MapKeysToLower returned %v, expected %v", lower, expected)
	}
}

func TestMapKeysToUpper(t *testing.T) {
	data := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": 2,
		},
	}
	upper := MapKeysToUpper(data)
	expected := map[string]interface{}{
		"A": 1,
		"B": map[string]interface{}{
			"C": 2,
		},
	}
	if !reflect.DeepEqual(upper, expected) {
		t.Errorf("MapKeysToUpper returned %v, expected %v", upper, expected)
	}
}

func TestConvertMapValues(t *testing.T) {
	data := map[string]interface{}{
		"a": 1,
		"b": "2",
	}
	converted := ConvertMapValues(data, func(v interface{}) interface{} {
		return fmt.Sprintf("%v", v)
	})
	expected := map[string]interface{}{
		"a": "1",
		"b": "2",
	}
	if !reflect.DeepEqual(converted, expected) {
		t.Errorf("ConvertMapValues returned %v, expected %v", converted, expected)
	}
}
