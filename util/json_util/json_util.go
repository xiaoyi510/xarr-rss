package json_util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyJson(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return []byte(""), err
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return []byte(""), err

	}

	return out.Bytes(), nil
}

func JsonMarshalToString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json 错误", v, err)
		return ""
	}
	return string(b)
}
