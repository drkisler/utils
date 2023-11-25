package utils

import (
	"fmt"
	"time"
)

type TMatchMapKey struct {
	Data map[string]string
}

func (ms *TMatchMapKey) GetInt(source map[string]interface{}, key string) (int, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0, nil
	}
	v, ok := result.(int)
	if !ok {
		return 0, fmt.Errorf("%s不是整数类型", key)
	}
	return v, nil
}

func (ms *TMatchMapKey) GetInt8(source map[string]interface{}, key string) (int8, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0, nil
	}
	v, ok := result.(int8)
	if !ok {
		return 0, fmt.Errorf("%s不是整数类型", key)
	}
	return v, nil
}

func (ms *TMatchMapKey) GetInt16(source map[string]interface{}, key string) (int16, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0, nil
	}
	v, ok := result.(int16)
	if !ok {
		return 0, fmt.Errorf("%s不是整数类型", key)
	}
	return v, nil
}
func (ms *TMatchMapKey) GetInt32(source map[string]interface{}, key string) (int32, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0, nil
	}
	switch v := result.(type) {
	case int:
		return int32(v), nil
	case int8:
		return int32(v), nil
	case int16:
		return int32(v), nil
	case int32:
		return v, nil
	case int64:
		return int32(v), nil
	default:
		return 0, fmt.Errorf("%s不是整数类型", key)
	}

}
func (ms *TMatchMapKey) GetInt64(source map[string]interface{}, key string) (int64, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0, nil
	}
	switch v := result.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil

	default:
		return 0, fmt.Errorf("%s不是整数类型", key)
	}

}
func (ms *TMatchMapKey) GetString(source map[string]interface{}, key string) (string, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return "", fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return "", nil
	}
	v, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("%s不是字符串类型", key)
	}
	return v, nil
}
func (ms *TMatchMapKey) GetFloat32(source map[string]interface{}, key string) (float32, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0.0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0.0, nil
	}
	switch v := result.(type) {
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	default:
		return 0.0, fmt.Errorf("%s不是浮点数类型", key)
	}
}
func (ms *TMatchMapKey) GetFloat64(source map[string]interface{}, key string) (float64, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return 0.0, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return 0.0, nil
	}

	//v := reflect.ValueOf(result)
	switch v := result.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0.0, fmt.Errorf("%s不是浮点数类型", key)
	}
}

func (ms *TMatchMapKey) GetBool(source map[string]interface{}, key string) (bool, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return false, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return false, nil
	}
	v, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("%s不是布尔类型", key)
	}
	return v, nil
}

func (ms *TMatchMapKey) GetTime(source map[string]interface{}, key string) (time.Time, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return time.Time{}, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return time.Time{}, nil
	}
	v, ok := result.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("%s不是时间类型", key)
	}
	return v, nil
}

func (ms *TMatchMapKey) GetSlice(source map[string]interface{}, key string) ([]interface{}, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return nil, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return nil, nil
	}
	v, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s不是切片类型", key)
	}
	return v, nil
}

func (ms *TMatchMapKey) GetStruct(source map[string]interface{}, key string) (interface{}, error) {
	result, ok := source[ms.Data[key]]
	if !ok {
		return nil, fmt.Errorf("%s不存在", key)
	}
	if result == nil {
		return nil, nil
	}
	_, ok = result.(struct{})
	if !ok {
		return nil, fmt.Errorf("%s不是结构体类型", key)
	}
	return result, nil
}
