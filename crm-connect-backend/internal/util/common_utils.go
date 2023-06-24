package util

import (
	"connect-crm-backend/crm-connect-backend/internal/constant"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var reservedResourceVerbs = []string{"assign", "expire"}

// Min ...
func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Generic contains which accepts any type of array ...
func Contains(arr interface{}, element interface{}) int {
	contentValue := reflect.ValueOf(arr)

	for i := 0; i < contentValue.Len(); i++ {
		if contentValue.Index(i).Interface() == element {
			return i
		}
	}

	return -1
}

// Split strings by splitBy and return the idx-th element
func SplitAndReturn(str string, splitBy string, idx int) string {
	splits := strings.Split(str, splitBy)
	if len(splits) <= idx {
		return constant.EmptyString
	}
	return splits[idx]
}

// GeneralizeURL ...
func GeneralizeURL(url string) string {
	urlSplits := strings.Split(strings.Trim(url, "/v1/"), "/")

	var generalizeURLBuilder strings.Builder

	generalizeURLBuilder.WriteString("/")
	// todo: use a better logic here
	for i := range urlSplits {
		if i == 3 {
			if containsReservedResourceVerb(urlSplits[i]) {
				generalizeURLBuilder.WriteString(urlSplits[i] + "/")
			} else {
				generalizeURLBuilder.WriteString(":id/")
			}
		} else {
			generalizeURLBuilder.WriteString(urlSplits[i] + "/")
		}
	}

	return generalizeURLBuilder.String()
}

// GetRequestID ...
func GetRequestID(r *http.Request) string {
	//if reqID := middleware.GetReqID(r.Context()); reqID != "" {
	//	return reqID
	//}

	return ""
}

func GetTimestampInMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Rot18 = Rot5 + Rot13
func Rot18(s string) string {
	var result []rune
	rot5map := map[rune]rune{
		'0': '5',
		'1': '6',
		'2': '7',
		'3': '8',
		'4': '9',
		'5': '0',
		'6': '1',
		'7': '2',
		'8': '3',
		'9': '4',
	}

	for _, i := range s {
		switch {
		case !unicode.IsLetter(i) && !unicode.IsNumber(i):
			result = append(result, i)
		case i >= 'A' && i <= 'Z':
			result = append(result, 'A'+(i-'A'+13)%26)
		case i >= 'a' && i <= 'z':
			result = append(result, 'a'+(i-'a'+13)%26)
		case i >= '0' && i <= '9':
			result = append(result, rot5map[i])
		case unicode.IsSpace(i):
			result = append(result, ' ')
		}
	}
	return fmt.Sprintf(string(result[:]))
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func ToString(v interface{}) (string, error) {
	vType := reflect.ValueOf(v).Kind()

	switch vType {
	case reflect.String:
		return v.(string), nil
	case reflect.Int:
		return strconv.Itoa(v.(int)), nil
	case reflect.Int64:
		return Int64ToString(v.(int64)), nil
	default:
		return constant.EmptyString, errors.New("Type not convertible to string")
	}
}

func StringToInt64(str string) int64 {
	num, _ := strconv.Atoi(str)
	return int64(num)
}

func ReadConfigFile(fileName string, v interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func NestedMapLookupWithArray(m map[string]interface{}, path []interface{}) (interface{}, error) {
	var ok bool
	var returnVal interface{}

	if len(path) == 0 {
		return nil, errors.New("invalid path")
	}

	switch x := path[0].(type) {
	case []interface{}:
		arrKey, ok := x[0].(string)
		if !ok {
			errMsg := fmt.Sprintf("malformed structure at %v", x[0])
			return nil, errors.New(errMsg)
		}

		index, ok := x[1].(int)
		if !ok {
			errMsg := fmt.Sprintf("malformed structure at %v", x[1])
			return nil, errors.New(errMsg)
		}
		arr, ok := m[arrKey].([]interface{})
		if !ok {
			errMsg := fmt.Sprintf("malformed structure at %v", x[0])
			return nil, errors.New(errMsg)
		}
		if index >= len(arr) {
			errMsg := fmt.Sprintf("array index %v out of range", index)
			return nil, errors.New(errMsg)
		}
		returnVal = arr[index]
	case string:
		returnVal, ok = m[x]
		if !ok {
			errMsg := fmt.Sprintf("key (%s) not found", path[0])
			return nil, errors.New(errMsg)
		}
	default:
		errMsg := fmt.Sprintf("malformed structure at %v", path[0])
		return nil, errors.New(errMsg)
	}

	if len(path) == 1 { // we've reached the final key
		return returnVal, nil
	}

	m, ok = returnVal.(map[string]interface{}) //go to next level in map

	if !ok { // final key not reached, map is at nesting depth
		errMsg := fmt.Sprintf("malformed structure at %v", path[0])
		return nil, errors.New(errMsg)
	}
	return NestedMapLookupWithArray(m, path[1:])
}

func ObjectFieldLookup(object interface{}, field string) (interface{}, error) {
	r := reflect.ValueOf(object)
	f := reflect.Indirect(r).FieldByName(field)

	switch f.Kind() {
	case reflect.Bool:
		return f.Bool(), nil
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return int(f.Int()), nil
	case reflect.String:
		return f.String(), nil
	case reflect.Float64:
		return f.Float(), nil
	default:
		errMsg := fmt.Sprintf("Unsupported type:(%s) for field:(%s)", r.Kind(), field)
		return nil, errors.New(errMsg)
	}

}

func GetEnvironmentVariable(varName string) (string, error) {
	value, ok := os.LookupEnv(varName)
	if !ok {
		errMsg := fmt.Sprintf("env variable %s not defined", varName)
		return constant.EmptyString, errors.New(errMsg)
	}
	return value, nil
}

func ConvertEpochTimeToDate(epochTime int64) string {
	date := time.Unix(epochTime, 0)
	return date.Format(constant.YYYYMMDD)
}

func containsReservedResourceVerb(input string) bool {
	for _, rv := range reservedResourceVerbs {
		if rv == input {
			return true
		}
	}
	return false
}
