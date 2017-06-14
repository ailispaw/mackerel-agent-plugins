package mpprometheus

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
)

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func PrintInJSON(out io.Writer, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	indented := new(bytes.Buffer)
	if err := json.Indent(indented, data, "", "  "); err != nil {
		return err
	}
	indented.WriteString("\n")

	if _, err := io.Copy(out, indented); err != nil {
		return err
	}

	return nil
}
