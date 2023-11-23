package utils

import (
	"fmt"
	"strings"
)

type TEnArrStr struct {
	ArrayString *[]string
}

func (arrStr *TEnArrStr) Load(source []string) {
	arrStr.ArrayString = &source
}

func (arrStr *TEnArrStr) MatchMapKey(source map[string]interface{}) map[string]string {
	var result = make(map[string]string)
	for _, str := range *(arrStr.ArrayString) {
	innerLoop:
		for key := range source {
			if strings.ToLower(key) == strings.ToLower(str) {
				result[str] = key
				break innerLoop
			}
		}
	}
	return result
}

func (arrStr *TEnArrStr) Exists(target string) bool {
	for _, str := range *arrStr.ArrayString {
		if str == target {
			return true
		}
	}
	return false
}

func (arrStr *TEnArrStr) IndexOf(target string, skip int) int {
	cnt := 0
	for iIndex, str := range *arrStr.ArrayString {
		if str == target {
			if cnt < skip {
				cnt++
				continue
			}
			return iIndex
		}
	}
	return -1
}

func (arrStr *TEnArrStr) Values(iIndex int) (string, error) {
	if iIndex < 0 || iIndex >= len(*arrStr.ArrayString) {
		return "", fmt.Errorf("out of index")
	}
	return (*arrStr.ArrayString)[iIndex], nil
}
