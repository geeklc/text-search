package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 根据数据行数和批次大小计算批次的数
func GetFloorBatchNum(rows int, batchSize int) int {
	if rows%batchSize > 0 {
		return (rows / batchSize) + 1
	} else {
		return rows / batchSize
	}
}

// 根据开始和结束标识截取数据 取值为[star, end)
func SliceRange[T any | int](slice []T, start, end int) []T {
	if slice == nil || len(slice) == 0 {
		return nil
	}
	var length = len(slice)
	if end > length {
		end = length
	}
	return slice[start:end]
}

// 根据时间格式转换
func GetNowTimeStr(format string) string {
	return time.Now().Format(format)
}

// 根据时间格式转换
func GetTimeStr(time time.Time, format string) string {
	if time.IsZero() {
		return ""
	}
	return time.Format(format)
}

// 拆分字符串，并为每个拆分的串去掉以空格开头的字符
func SplitAndTrimBegin(s, sep string) []string {
	infos := strings.Split(s, sep)
	for i := 0; i < len(infos); i++ {
		infos[i] = strings.TrimPrefix(infos[i], " ")
	}
	return infos
}

func AnyToMysqlFiled(i interface{}) string {
	if i == nil {
		return "NULL"
	}
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "NULL"
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:

		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64:
		return fmt.Sprintf("(%g+%gi)", real(v.Complex()), imag(v.Complex()))
	case reflect.Complex128:
		return fmt.Sprintf("(%g+%gi)", real(v.Complex()), imag(v.Complex()))
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Slice, reflect.Map, reflect.Struct, reflect.Array:
		str, _ := json.Marshal(i)
		return string(str)
	default:
		return ""
	}
}

func Md5(data string) string {
	// 创建一个 MD5 哈希对象
	hash := md5.New()
	// 写入数据
	hash.Write([]byte(data))
	// 计算哈希值
	hashBytes := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}

// 判断数组包含某个值
func Contains[T comparable](slice []T, value T) bool {
	if slice == nil || len(slice) == 0 {
		return false
	}
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GetMapValues(infos map[string]interface{}) []interface{} {
	if infos == nil {
		return make([]any, 0)
	}
	args := make([]interface{}, len(infos))
	n := 0
	for _, value := range infos {
		args[n] = value
		n++
	}
	return args
}

func InterfaceArrayCatMap(infos []interface{}) ([]map[string]interface{}, error) {
	if infos == nil || len(infos) == 0 {
		return nil, nil
	}
	// 遍历切片，并将每个元素转换为 map
	var maps []map[string]interface{}
	for _, item := range infos {
		mapItem, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Expected a map in the slice")
			return nil, errors.New("Expected a map in the slice")
		}
		maps = append(maps, mapItem)
	}
	return maps, nil
}
