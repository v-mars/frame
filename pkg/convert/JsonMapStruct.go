package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
)


func StructToMapOut(m interface{}, out *map[string]interface{}) error {
	tmp, err := json.Marshal(m)
	if err!=nil{
		return err
	}
	err = json.Unmarshal(tmp, &out)
	if err!=nil{
		return err
	}
	return nil
}

func StructToMap(m interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	tmp, err := json.Marshal(m)
	if err!=nil{
		return nil,err
	}
	err = json.Unmarshal(tmp, &out)
	if err!=nil{
		return nil,err
	}
	return out, nil
}

func StructToMapSlice(m interface{}) ([]map[string]interface{}, error) {
	var out []map[string]interface{}
	tmp, err := json.Marshal(m)
	if err!=nil{
		return nil,err
	}
	err = json.Unmarshal(tmp, &out)
	if err!=nil{
		return nil,err
	}
	return out, nil
}

func StructToMapViaReflect(m interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	elem := reflect.ValueOf(&m).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		r[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return r
}

// StructToMapByReflect ToMap 结构体转为Map[string]interface{}
func StructToMapByReflect(in interface{}, tagName string) (map[string]interface{}, error){
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {  // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	relType := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := relType.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		} else {
			out[relType.Field(i).Name] = v.Interface()
		}
	}
	return out, nil
}