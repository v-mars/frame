package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/vmihailenco/msgpack"
)

func AnyToAny(src interface{}, des interface{}) error {
	bytesData, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytesData, des)
	return err
}
func AnyToAnyV2(src interface{}, des interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(src)
	err = json.NewDecoder(&buf).Decode(&des)
	return err
}

func GobAnyToAny(src interface{}, des interface{}) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(src)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	//dec := gob.NewDecoder(&buf)
	err = decoder.Decode(des)
	return err
}

func MsgpackAnyToAny(src interface{}, des interface{}) error {
	bytesData, err := msgpack.Marshal(src)
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal(bytesData, des)
	return err
}
