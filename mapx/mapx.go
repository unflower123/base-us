package mapx

import (
	"encoding/json"
	"strings"
)

// MapToJson convert map to JSON string
func MapToJson[K comparable, V any](data map[K]V) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JsonToMap convert JSON string to map
func JsonToMap[K comparable, V any](jsonStr string) (map[K]V, error) {
	var data map[K]V
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// StructToMap convert the struct to map
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// MapToStruct convert map to struct
func MapToStruct(data map[string]interface{}, obj interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, obj)
	if err != nil {
		return err
	}

	return nil
}

// MergeMaps deep merge two maps
func MergeMaps[K comparable, V any](map1, map2 map[K]V, isOverwrite bool) map[K]V {
	mergedMap := make(map[K]V)
	// copy map1
	for k, v := range map1 {
		mergedMap[k] = v
	}

	// recursively merge the contents of map2
	for k, v := range map2 {
		if existing, ok := mergedMap[k]; ok {
			// if the key already exists and the value is map, recursively merge
			if existingMap, ok := any(existing).(map[K]V); ok {
				if newMap, ok := any(v).(map[K]V); ok {
					mergedMap[k] = any(MergeMaps[K, V](existingMap, newMap, isOverwrite)).(V)
					continue
				}
			}
			if !isOverwrite {
				continue
			}
		}
		mergedMap[k] = v
	}

	return mergedMap
}

// FilterMap  filter key value pairs in the map
func FilterMap[K comparable, V any](data map[K]V, filterFunc func(key K, value V) bool) map[K]V {
	filteredMap := make(map[K]V)
	for k, v := range data {
		if filterFunc(k, v) {
			filteredMap[k] = v
		}
	}
	return filteredMap
}

// MapValues perform mapping transformation on each value in the map
func MapValues[K comparable, V any, R any](data map[K]V, mapFunc func(key K, value V) R) map[K]R {
	mappedMap := make(map[K]R)
	for k, v := range data {
		mappedMap[k] = mapFunc(k, v)
	}
	return mappedMap
}

// MapKeysToLower convert the keys of the map to lower (support nested maps)
func MapKeysToLower(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		lowerKey := strings.ToLower(k)
		if nestedMap, ok := v.(map[string]interface{}); ok {
			newMap[lowerKey] = MapKeysToLower(nestedMap)
		} else {
			newMap[lowerKey] = v
		}
	}
	return newMap
}

// MapKeysToUpper convert the keys of the map to upper (support nested maps)
func MapKeysToUpper(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		lowerKey := strings.ToUpper(k)
		if nestedMap, ok := v.(map[string]interface{}); ok {
			newMap[lowerKey] = MapKeysToUpper(nestedMap)
		} else {
			newMap[lowerKey] = v
		}
	}
	return newMap
}

// ConvertMapValues convert the values in the map to the specified type
func ConvertMapValues(data map[string]interface{}, convertFunc func(interface{}) interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for k, v := range data {
		converted[k] = convertFunc(v)
	}
	return converted
}
