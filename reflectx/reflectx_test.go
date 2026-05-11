package reflectx

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email"`
	Tags  []string `json:"tags"`
}

func TestGetInterfaceValList(t *testing.T) {
	user := TestStruct{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
		Tags:  []string{"admin", "user"},
	}

	// 测试获取符合标签条件的字段值
	values := GetInterfaceValList(user, "json", "name")
	if len(values) != 1 {
		t.Errorf("GetInterfaceValList returned %d values, expected 1", len(values))
	}
	if values[0].Interface() != "Alice" {
		t.Errorf("GetInterfaceValList returned %v, expected 'Alice'", values[0].Interface())
	}

	// 测试获取嵌套字段值
	type NestedStruct struct {
		User TestStruct `json:"user"`
	}
	nested := NestedStruct{User: user}
	values = GetInterfaceValList(nested, "json", "name")
	if len(values) != 1 {
		t.Errorf("GetInterfaceValList returned %d values, expected 1", len(values))
	}
	if values[0].Interface() != "Alice" {
		t.Errorf("GetInterfaceValList returned %v, expected 'Alice'", values[0].Interface())
	}
}

func TestOperateInterface(t *testing.T) {
	user := TestStruct{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
		Tags:  []string{"admin", "user"},
	}

	// 测试对字段执行操作
	count := 0
	OperateInterface(user, func(rv reflect.Value, rt reflect.StructField) {
		count++
	})
	if count != 4 {
		t.Errorf("OperateInterface processed %d fields, expected 4", count)
	}
}

func TestGetList(t *testing.T) {
	user := TestStruct{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
		Tags:  []string{"admin", "user", "admin"}, // 包含重复值
	}

	// 测试获取字符串切片字段值
	tags := GetList[string](user, "json", "tags")
	expectedTags := []string{"admin", "user"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("GetList returned %v, expected %v", tags, expectedTags)
	}

	// 测试获取整型字段值
	ages := GetList[int](user, "json", "age")
	expectedAges := []int{25}
	if !reflect.DeepEqual(ages, expectedAges) {
		t.Errorf("GetList returned %v, expected %v", ages, expectedAges)
	}
}

func TestIsPublicField(t *testing.T) {
	type TestStruct struct {
		Name  string
		age   int
		Email string
	}

	rt := reflect.TypeOf(TestStruct{})
	if !isPublicField(rt.Field(0)) {
		t.Error("isPublicField returned false for public field, expected true")
	}
	if isPublicField(rt.Field(1)) {
		t.Error("isPublicField returned true for private field, expected false")
	}
}

func TestStringContains(t *testing.T) {
	if !stringContains("a", "a", "b", "c") {
		t.Error("stringContains returned false, expected true")
	}
	if stringContains("d", "a", "b", "c") {
		t.Error("stringContains returned true, expected false")
	}
}
