package rest

import (
	"encoding/json"
)

type Response struct {
	Code int32 `json:"code"`
	Desc string `json:"desc"`
}

func (r *Response) Byte() ([]byte,error) {
	return json.Marshal(r)
}

