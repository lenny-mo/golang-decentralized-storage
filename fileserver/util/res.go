// this file is to declare a data type and some func to convert info to json format
package util

import (
	"encoding/json"
	"log"
)

type ResponseMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSON2Bytes convert json to bytes
func (obj *ResponseMessage) JSON2Bytes() []byte {
	res, err := json.Marshal(obj)
	if err != nil {
		log.Println("json marshal err: ", err.Error())
		return nil
	}
	return res
}

// JSON2String convert json to string
func (obj *ResponseMessage) JSON2String() string {
	res := obj.JSON2Bytes()
	return string(res)
}
