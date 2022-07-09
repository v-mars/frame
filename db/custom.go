package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	bytesUnit "github.com/docker/go-units"
	"github.com/v-mars/frame/pkg/convert"
)

type IntArray []int

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *IntArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %s", value)
	}
	var err error
	result := IntArray{}
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &result)
	}
	*j = result
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j IntArray) Value() (driver.Value, error) {
	if &j == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(j)
	return string(bytes), err
}

type StringArray []string

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *StringArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %s", value)
	}
	var err error
	result := StringArray{}
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &result)
	}
	*j = result
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j StringArray) Value() (driver.Value, error) {
	if &j == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(j)
	return string(bytes), err
}

type IntNestArray [][]int

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *IntNestArray) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %s", value)
	}
	var err error
	result := IntNestArray{}
	if len(bytes) != 0 {
		err = json.Unmarshal(bytes, &result)
	}
	*j = result
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j IntNestArray) Value() (driver.Value, error) {
	if &j == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(j)
	return string(bytes), err
}

type BytesUnit string

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *BytesUnit) Scan(value interface{}) error {
	var bytes float64
	bytes, err := convert.ToFloat64E(value)
	if err != nil {
		return fmt.Errorf("failed to convert.ToFloat64: %s", value)
	}

	*j = "0"
	if bytes != 0 {
		*j = BytesUnit(bytesUnit.BytesSize(bytes))
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j BytesUnit) Value() (driver.Value, error) {
	if &j == nil || len(string(j)) == 0 {
		return nil, nil
	}
	value, err := bytesUnit.RAMInBytes(string(j))
	return value, err
}
