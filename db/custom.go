package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	bytesUnit "github.com/docker/go-units"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/utils"
	"log"
	"time"
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
	*j = "0"
	var bytes float64
	bytes = convert.ToFloat64(value)
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

type Bool bool

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *Bool) Scan(value interface{}) error {
	*j = Bool(false)
	if value != nil {
		*j = Bool(convert.ToBool(value))
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j Bool) Value() (driver.Value, error) {
	if &j == nil {
		return nil, nil
	}
	value := convert.ToString(j)
	return value, nil
}

type Int int

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *Int) Scan(value interface{}) error {
	*j = 0
	if value != nil {
		*j = Int(convert.ToInt(value))
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j Int) Value() (driver.Value, error) {
	if &j == nil {
		return nil, nil
	}
	value := convert.ToString(j)
	return value, nil
}

type MillisTime string

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
//func (j MillisTime) MarshalJSON() ([]byte, error) {
//	return []byte(strconv.Quote(string(j))), nil
//}

// UnmarshalJSON 反序列化
func (j *MillisTime) UnmarshalJSON(data []byte) error {
	*j = MillisTime(data)
	return nil
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *MillisTime) Scan(value interface{}) error {
	*j = "0"
	if value != nil {
		*j = MillisTime((time.Millisecond * time.Duration(convert.ToInt64(value))).String())
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j MillisTime) Value() (driver.Value, error) {
	if len(j) == 0 || j == "0" || j == "0s" {
		return "0", nil
	}
	duration, err := time.ParseDuration(string(j))
	if err != nil {
		log.Printf("MillisTime time.ParseDuration failed: %s", err)
		return "0", nil
	}
	return convert.ToString(duration.Milliseconds()), nil
}

type SecTime string

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
//func (j SecTime) MarshalJSON() ([]byte, error) {
//
//	//formatted := ""
//	//if len(j.Value()) > 0 && j.String() != "null" {
//	//	formatted = fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
//	//}
//	return []byte(strconv.Quote(string(j))), nil
//}

// UnmarshalJSON 反序列化
func (j *SecTime) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	*j = SecTime((time.Millisecond * time.Duration(convert.ToFloat64(string(data))*1000)).String())
	return nil
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *SecTime) Scan(value interface{}) error {
	*j = "0s"
	if value != nil {
		*j = SecTime((time.Millisecond * time.Duration(convert.ToFloat64(value)*1000)).String())
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j SecTime) Value() (driver.Value, error) {
	if len(j) == 0 || j == "0" || j == "0s" {
		return "0", nil
	}
	duration, err := time.ParseDuration(string(j))
	if err != nil {
		log.Printf("SecTime time.ParseDuration failed: %s", err)
		return "0", nil
	}
	return convert.ToString(duration.Seconds()), nil
}

type AesStr string

var DefaultAesKey = "AY3b5Z78806GorMa"

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *AesStr) Scan(value interface{}) error {
	//*j = ""
	var str = convert.ToString(value)
	if len(str) > 0 {
		*j = AesStr(utils.DeTxtByAes(str, DefaultAesKey))
	}
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j AesStr) Value() (driver.Value, error) {
	if len(j) == 0 || len(string(j)) == 0 {
		return "", nil
	}
	return utils.EnTxtByAes(string(j), DefaultAesKey), nil
}
