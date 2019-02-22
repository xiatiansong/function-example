package simplejson

import (
	"reflect"
	"fmt"
	"math"
	"errors"
)

import (
	"encoding/json"
	"strconv"
)

func MustString(v interface{}, defaultValue string) string {
	switch tv := v.(type) {
	case string:
		return tv
	case []byte:
		return string(tv)
	case int64:
		return strconv.FormatInt(int64(tv), 10)
	case uint64:
		return strconv.FormatUint(uint64(tv), 10)
	case int32:
		return strconv.FormatInt(int64(tv), 10)
	case uint32:
		return strconv.FormatUint(uint64(tv), 10)
	case int:
		return strconv.Itoa(int(tv))
	case int16:
		return strconv.FormatInt(int64(tv), 10)
	case uint16:
		return strconv.FormatUint(uint64(tv), 10)
	case float32:
		return strconv.FormatFloat(float64(tv), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(tv, 'f', -1, 64)
	case json.Number:
		return tv.String()
	case bool:
		if tv {
			return "true"
		} else {
			return "false"
		}
	}
	return defaultValue
}

func TryString(v interface{}) (string, bool) {
	switch tv := v.(type) {
	case string:
		return tv, true
	case []byte:
		return string(tv), true
	case int64:
		return strconv.FormatInt(int64(tv), 10), true
	case uint64:
		return strconv.FormatUint(uint64(tv), 10), true
	case int32:
		return strconv.FormatInt(int64(tv), 10), true
	case uint32:
		return strconv.FormatUint(uint64(tv), 10), true
	case int16:
		return strconv.FormatInt(int64(tv), 10), true
	case uint16:
		return strconv.FormatUint(uint64(tv), 10), true
	case float32:
		return strconv.FormatFloat(float64(tv), 'f', -1, 64), true
	case float64:
		return strconv.FormatFloat(float64(tv), 'f', -1, 64), true
	case int:
		return strconv.Itoa(int(tv)), true
	case json.Number:
		return tv.String(), true
	case bool:
		if tv {
			return "true", true
		} else {
			return "false", true
		}
	}
	return "", false
}

func MustInt64(v interface{}, defaultValue int64) int64 {
	if v == nil {
		return defaultValue
	}
	switch tv := v.(type) {
	case []byte:
		res, err := strconv.ParseInt(string(tv), 10, 0)
		if err != nil {
			return defaultValue
		}
		return res
	case string:
		res, err := strconv.ParseInt(tv, 10, 0)
		if err != nil {
			return defaultValue
		}
		return res
	case int64:
		return tv
	case uint64:
		if tv > uint64(math.MaxInt64) {
			return defaultValue
		}
		return int64(tv)
	case int32:
		return int64(tv)
	case uint32:
		return int64(tv)
	case int:
		return int64(tv)
	case float32:
		if tv > float32(math.MaxInt64) {
			return defaultValue
		}
		return int64(tv)
	case float64:
		if tv > float64(math.MaxInt64) {
			return defaultValue
		}
		return int64(tv)
	case json.Number:
		val, err := tv.Int64()
		if err == nil {
			return val
		}
	}
	return defaultValue
}

func MustFloat64(v interface{}, defaultValue float64) float64 {
	switch tv := v.(type) {
	case []byte:
		res, err := strconv.ParseFloat(string(tv), 0)
		if err != nil {
			return defaultValue
		}
		return res
	case string:
		res, err := strconv.ParseFloat(tv, 0)
		if err != nil {
			return defaultValue
		}
		return res
	case int64:
		return float64(tv)
	case uint64:
		return float64(tv)
	case int32:
		return float64(tv)
	case uint32:
		return float64(tv)
	case int:
		return float64(tv)
	case float32:
		return float64(tv)
	case float64:
		return tv
	case json.Number:
		val, err := tv.Float64()
		if err == nil {
			return val
		}
	}
	return defaultValue
}

func TryFloat64(v interface{}) (float64, bool) {
	switch tv := v.(type) {
	case []byte:
		res, err := strconv.ParseFloat(string(tv), 0)
		if err != nil {
			return 0, false
		}
		return res, true
	case string:
		res, err := strconv.ParseFloat(tv, 0)
		if err != nil {
			return 0, false
		}
		return res, true
	case int64:
		return float64(tv), true
	case uint64:
		return float64(tv), true
	case int32:
		return float64(tv), true
	case uint32:
		return float64(tv), true
	case int:
		return float64(tv), true
	case float32:
		return float64(tv), true
	case float64:
		return tv, true
	case json.Number:
		val, err := tv.Float64()
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

// This is a safe convert since string can always convert to interface{}
func StringStringMap2StringInterfaceMap(input map[string]string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range input {
		res[k] = v
	}
	return res
}

func ConvertToInt64(v interface{}) (int64, error) {
	defaultValue := int64(-1)
	if v == nil {
		return -1, errors.New("input is nil")
	}
	switch tv := v.(type) {
	case []byte:
		res, err := strconv.ParseInt(string(tv), 10, 0)
		if err != nil {
			return defaultValue, err
		}
		return res, nil
	case string:
		res, err := strconv.ParseInt(tv, 10, 0)
		if err != nil {
			return defaultValue, err
		}
		return res, nil
	case int64:
		return tv, nil
	case uint64:
		if tv > uint64(math.MaxInt64) {
			return defaultValue, errors.New("input number out of range")
		}
		return int64(tv), nil
	case int32:
		return int64(tv), nil
	case uint32:
		return int64(tv), nil
	case int:
		return int64(tv), nil
	case float32:
		if tv > float32(math.MaxInt64) {
			return defaultValue, errors.New("input number out of range")
		}
		return int64(tv), nil
	case float64:
		if tv > float64(math.MaxInt64) {
			return defaultValue, errors.New("input number out of range")
		}
		return int64(tv), nil
	case json.Number:
		val, err := tv.Int64()
		if err == nil {
			return val, nil
		}
		return defaultValue, err
	}
	return defaultValue, fmt.Errorf("input number type err type=%v", reflect.TypeOf(v))
}


func MustStringArray(v interface{}, defaultValue []string) []string {
	var ret []string
	x := reflect.ValueOf(v)
	switch x.Kind() {
	case reflect.Array, reflect.String, reflect.Slice:
		for i := 0; i < x.Len(); i++ {
			val := x.Index(i).Interface()
			if val == nil {
				return defaultValue
			}
			if valStr, ok := TryString(val); ok {
				ret = append(ret, valStr)
			} else {
				return defaultValue
			}
		}
		return ret
	default:
		return defaultValue
	}
}

func TryBool(v interface{}) (ret bool, isbool bool) {
	if v == nil {
		return false, false
	}
	switch tv := v.(type) {
	case bool:
		return tv, true
	case string:
		//Attention:
		//   strconv.ParseBool() think "0","t","T","1","f","F" as bool,
		//   but it may not used as bool for outer func.
		//	 So only return value that must be used as bool.
		switch tv {
		case "true", "TRUE", "True":
			return true, true
		case "false", "FALSE", "False":
			return false, true
		default:
			return false, false
		}
	default:
		return false, false
	}
}
