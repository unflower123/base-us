package slicex

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFilterSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	filtered := FilterSlice(slice, func(item int) bool {
		return item > 2
	})
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("FilterSlice returned %v, expected %v", filtered, expected)
	}
}

func TestUniqueSlice(t *testing.T) {
	slice := []int{1, 2, 2, 3, 3, 3}
	unique := UniqueSlice(slice)
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(unique, expected) {
		t.Errorf("UniqueSlice returned %v, expected %v", unique, expected)
	}
}

func TestRemoveEmptySlice(t *testing.T) {
	slice := []string{"a", "", "b", "", "c"}
	cleaned := RemoveEmptySlice(slice)
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(cleaned, expected) {
		t.Errorf("RemoveEmptySlice returned %v, expected %v", cleaned, expected)
	}
}

func TestSortSlice(t *testing.T) {
	slice := []int{3, 1, 2}
	sorted := SortSlice(slice, func(a, b int) bool {
		return a < b
	})
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("SortSlice returned %v, expected %v", sorted, expected)
	}
}

func TestReverseSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	reversed := ReverseSlice(slice)
	expected := []int{3, 2, 1}
	if !reflect.DeepEqual(reversed, expected) {
		t.Errorf("ReverseSlice returned %v, expected %v", reversed, expected)
	}
}

func TestAppendIfNotExists(t *testing.T) {
	slice := []int{1, 2, 3}
	slice = AppendIfNotExists(slice, 4)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("AppendIfNotExists returned %v, expected %v", slice, expected)
	}

	slice = AppendIfNotExists(slice, 2)
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("AppendIfNotExists returned %v, expected %v", slice, expected)
	}
}

func TestJoin(t *testing.T) {
	strs := []string{"a", "b", "c"}
	joined := Join(strs, ", ")
	expected := "a, b, c"
	if joined != expected {
		t.Errorf("Join returned %v, expected %v", joined, expected)
	}
}

func TestIntersectSlice(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{2, 3, 4}
	intersect := IntersectSlice(slice1, slice2)
	expected := []int{2, 3}
	if !reflect.DeepEqual(intersect, expected) {
		t.Errorf("IntersectSlice returned %v, expected %v", intersect, expected)
	}
}

func TestDifferenceSlice(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{2, 3, 4}
	diff := DifferenceSlice(slice1, slice2)
	expected := []int{1}
	if !reflect.DeepEqual(diff, expected) {
		t.Errorf("DifferenceSlice returned %v, expected %v", diff, expected)
	}
}

func TestUnionSlice(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{2, 3, 4}
	union := UnionSlice(slice1, slice2)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(union, expected) {
		t.Errorf("UnionSlice returned %v, expected %v", union, expected)
	}
}

func TestConvertSliceElements(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := ConvertSliceElements(slice1, func(t int) string {
		return fmt.Sprintf("%d", t)
	})
	expected := []string{"1", "2", "3"}
	if !reflect.DeepEqual(slice2, expected) {
		t.Errorf("ConvertSliceElements returned %v, expected %v", slice2, expected)
	}
}
