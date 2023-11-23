package utils

import "fmt"

type TMatchMapKey struct {
	Data map[string]string
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
