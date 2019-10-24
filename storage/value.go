package storage

import (
	"fmt"
	"strconv"
)

const emptyFlag = 0x1E

// Value is the target string
type Value string

// Set string
func (s *Value) Set(v string) {
	if v != "" {
		*s = Value(v)
	} else {
		s.Clear()
	}
}

// Clear string
func (s *Value) Clear() {
	*s = Value(emptyFlag)
}

// Exist check string exist
func (s Value) Exist() bool {
	return string(s) != string(emptyFlag)
}

// Bool string to bool
func (s Value) Bool() (bool, error) {
	return strconv.ParseBool(s.String())
}

// Float64 string to float64
func (s Value) Float64() (float64, error) {
	return strconv.ParseFloat(s.String(), 64)
}

// Int64 string to int64
func (s Value) Int64() (int64, error) {
	v, err := strconv.ParseInt(s.String(), 10, 64)
	return v, err
}

// Uint string to uint
func (s Value) Uint() (uint, error) {
	v, err := strconv.ParseUint(s.String(), 10, 32)
	return uint(v), err
}

// Uint64 string to uint64
func (s Value) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(s.String(), 10, 64)
	return v, err
}

// String string to string
func (s Value) String() string {
	if s.Exist() {
		return string(s)
	}
	return ""
}
func (s Value) Binary() []byte {
	if s.Exist() {
		return []byte(string(s))
	}
	return nil
}

// ToValue interface to string
func ToValue(value interface{}) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f',10, 32)
	case float64:
		s = strconv.FormatFloat(v, 'f', 10, 64)
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int8:
		s = strconv.FormatInt(int64(v), 10)
	case int16:
		s = strconv.FormatInt(int64(v), 10)
	case int32:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(v, 10)
	case uint:
		s = strconv.FormatUint(uint64(v), 10)
	case uint8:
		s = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s = strconv.FormatUint(v, 10)
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

