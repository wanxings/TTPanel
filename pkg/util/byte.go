package util

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
)

func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
func FormatSize(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

// StructToJsonStr 结构体转json字符串
func StructToJsonStr(s interface{}) (string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	b, err := json.Marshal(v.Interface())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// JsonStrToStruct json字符串转结构体
func JsonStrToStruct(jsonStr string, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("must be a pointer")
	}
	if v.IsNil() {
		return fmt.Errorf("must not be nil")
	}
	err := json.Unmarshal([]byte(jsonStr), s)
	if err != nil {
		return err
	}
	return nil
}
