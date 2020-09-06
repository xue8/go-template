package gotemplate

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDefault(t *testing.T) {
//	reg := regexp.MustCompile(`(?=}})`)
//	s := `{{   dd }}
//{{- ddwdw -}}
//"{{- dsd -}}"`
//	aaa := reg.FindAllString(s, -1)
//	fmt.Println(aaa)


	ma := make(map[string]interface{})
	ma["a"] = "1111\nwdwde;%$()*"
	ma["b"] = 2222
	xxx, _ := json.Marshal(struct {
		A string
		B int
		C int
	}{"hello", 2, 3})
	ma["c"] = string(xxx)
	ma["d"] = []string{"xx", "xxx", "wwwwwww\n3"}
	t1 := Default()
	aaaaa, err := t1.ParseJSON(`{
  "sites": {
    "a": "{{ a}}",
    "b": "{{ b }}"
  },
  "c": "{{ c }}",
  "d": "{{- d -}}"
}`, ma)
	if err != nil {

	}
	fmt.Println(aaaaa)
}

