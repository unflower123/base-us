package reflectx

import (
	"base/slicex"
	"reflect"
	"strings"
	"unicode"
)

// GetInterfaceValList Recursive traversal of structures, pointers, and slices,
// Obtain field values that meet the label criteria
func GetInterfaceValList(v interface{}, relateKey, relate string) (r []reflect.Value) {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	switch rt.Kind() {
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			itemFieldV := rv.Field(i)
			itemFieldT := rt.Field(i)
			if !isPublicField(rt.Field(i)) {
				continue
			}
			tag := itemFieldT.Tag.Get(relateKey)
			if tag == "-" {
				continue
			}
			tags := strings.Split(tag, ",")
			if stringContains(relate, tags...) {
				r = append(r, itemFieldV)
			} else {
				itemV := itemFieldV.Interface()
				r = append(r, GetInterfaceValList(itemV, relateKey, relate)...)
			}
		}
	case reflect.Ptr:
		if rv.IsNil() {
			return
		}
		v = rv.Elem().Interface()
		r = append(r, GetInterfaceValList(v, relateKey, relate)...)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			itemV := rv.Index(i).Interface()
			r = append(r, GetInterfaceValList(itemV, relateKey, relate)...)
		}
	default:
		return
	}
	return
}

// OperateInterface Recursive traversal of structures, pointers, and slices,
// And perform custom operations
func OperateInterface(v interface{}, fn func(rv reflect.Value, rt reflect.StructField)) {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			if !isPublicField(rt.Field(i)) {
				continue
			}
			fn(rv.Field(i), rt.Field(i))
		}
	case reflect.Ptr:
		if rv.IsNil() {
			return
		}
		v = rv.Elem().Interface()
		OperateInterface(v, fn)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			itemV := rv.Index(i).Interface()
			OperateInterface(itemV, fn)
		}
	default:
		return
	}
	return
}

// GetList Recursive traversal of structures, pointers, and slices,
// extracting field values that meet specified label conditions,
// and converting them into slices of the specified type
func GetList[T any](v interface{}, relateKey, relate string) (r []T) {
	list := GetInterfaceValList(v, relateKey, relate)
	for _, t := range list {
		if t.Kind() == reflect.Slice {
			for i := 0; i < t.Len(); i++ {
				if val, ok := t.Index(i).Interface().(T); ok {
					r = append(r, val)
				}
			}
		} else if val, ok := t.Interface().(T); ok {
			r = append(r, val)
		}
	}
	r = slicex.FilterSlice(r, func(item T) bool {
		return !reflect.DeepEqual(item, reflect.Zero(reflect.TypeOf(item)).Interface())
	})
	r = slicex.UniqueSlice(r)
	return
}

func isPublicField(f reflect.StructField) bool {
	if len(f.Name) == 0 {
		return false
	}
	return unicode.IsUpper(rune(f.Name[0]))
}

func stringContains(v string, list ...string) bool {
	for _, i := range list {
		if i == v {
			return true
		}
	}
	return false
}
