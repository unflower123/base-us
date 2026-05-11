package slicex

import (
	"sort"
	"strings"
)

// FilterSlice filter slices and retain elements that meet the conditions
func FilterSlice[T any](slice []T, filterFunc func(item T) bool) []T {
	var result []T
	for _, item := range slice {
		if filterFunc(item) {
			result = append(result, item)
		}
	}
	return result
}

// UniqueSlice remove duplicate elements from the slice
func UniqueSlice[T any](slice []T) []T {
	seen := make(map[any]bool)
	var result []T
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// RemoveEmptySlice remove empty values from slices
func RemoveEmptySlice[T comparable](slice []T) []T {
	var result []T
	for _, v := range slice {
		var zero T
		if v != zero {
			result = append(result, v)
		}
	}
	return result
}

// SortSlice sort the slices
func SortSlice[T any](slice []T, lessFunc func(a, b T) bool) []T {
	sorted := make([]T, len(slice))
	copy(sorted, slice)
	sort.Slice(sorted, func(i, j int) bool {
		return lessFunc(sorted[i], sorted[j])
	})
	return sorted
}

// ReverseSlice reverse the order of slices
func ReverseSlice[T any](slice []T) []T {
	reversed := make([]T, len(slice))
	for i := 0; i < len(slice); i++ {
		reversed[len(slice)-1-i] = slice[i]
	}
	return reversed
}

// AppendIfNotExists add elements to the slice, but only when the element does not exist in the slice
func AppendIfNotExists[T comparable](slice []T, item T) []T {
	for _, v := range slice {
		if v == item {
			return slice
		}
	}
	return append(slice, item)
}

// Join concatenate the elements in the slice into a string using the specified delimiter
func Join(slice []string, sep string) string {
	strs := make([]string, len(slice))
	for i, item := range slice {
		strs[i] = item
	}
	return strings.Join(strs, sep)
}

// IntersectSlice obtain the intersection of two slices
func IntersectSlice[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	var result []T
	for _, v := range slice2 {
		if set[v] {
			result = append(result, v)
		}
	}
	return result
}

// DifferenceSlice obtain the difference set between two slices
func DifferenceSlice[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	for _, v := range slice2 {
		set[v] = true
	}
	var result []T
	for _, v := range slice1 {
		if !set[v] {
			result = append(result, v)
		}
	}
	return result
}

// UnionSlice obtain the union of two slices
func UnionSlice[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	var result []T
	for _, v := range slice1 {
		if !set[v] {
			set[v] = true
			result = append(result, v)
		}
	}
	for _, v := range slice2 {
		if !set[v] {
			set[v] = true
			result = append(result, v)
		}
	}
	return result
}

func ConvertSliceElements[T any, U any](input []T, converter func(T) U) []U {
	output := make([]U, len(input))
	for i, v := range input {
		output[i] = converter(v)
	}
	return output
}

func ConvertSliceElementsWithError[T any, U any](input []T, converter func(T) (u U, err error)) ([]U, error) {
	var err error
	output := make([]U, len(input))
	for i, v := range input {
		output[i], err = converter(v)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}
