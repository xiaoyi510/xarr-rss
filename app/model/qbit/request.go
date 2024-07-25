package qbit

import (
	"fmt"
	"reflect"
)

type QbitRequest struct {
	Uri    string
	Method string
	Data   string
}

func (r *QbitRequest) GetApiType() string {
	field, _ := reflect.TypeOf(*r).FieldByName("apiType")
	return field.Tag.Get("api_type")
}

func (r *QbitRequest) GetApiMethod() string {
	field, _ := reflect.TypeOf(*r).FieldByName("apiType")
	return field.Tag.Get("api_method")
}

func (r *QbitRequest) Gen(host string) interface{} {
	field, _ := reflect.TypeOf(*r).FieldByName("apiType")
	r.Uri = host + r.parseUri()
	r.Method = field.Tag.Get("req_method")
	return r
}

func (r *QbitRequest) parseUri() string {
	return fmt.Sprintf("/api/v2/%s/%s", r.GetApiType(), r.GetApiMethod())
}
