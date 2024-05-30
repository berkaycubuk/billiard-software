package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func FormatJSON(data []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", " ")

	if err != nil {
		fmt.Println(err.Error())
	}

	d := out.Bytes()
	return string(d)
}
