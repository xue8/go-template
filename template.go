package gotemplate

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	texttemplate "text/template"
)

type Template struct {
	template *texttemplate.Template
	left string
	right string
}

func Default() *Template {
	return New("{{", "}}", true)
}

func New(left, right string, require bool) *Template {
	t := texttemplate.New("default")
	if 	t.Delims(left, right); require {
		t.Option("missingkey=error")
	}
	return &Template{
		template:texttemplate.New("default"),
		left:left,
		right:right,
	}
}

func (t *Template) Parse(text string, data interface{}) ([]byte, error) {
	var err error
	if text, err = t.trim(text); err != nil {
		return nil, err
	}
	if 	_, err = t.template.Parse(text); err != nil {
		return nil, err
	}
	bf := &bytes.Buffer{}
	if err = t.template.Execute(bf, data); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (t *Template) ParseJSON(text string, data interface{}) ([]byte, error) {
	var err error
	if 	text, err = t.trim(text); err != nil {
		return nil, err
	}
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			var vBytes []byte
			if vBytes, err = json.Marshal(v); err != nil {
				return nil, err
			}
			m[k] = string(vBytes)
		}
	}
	if 	_, err = t.template.Parse(text); err != nil {
		return nil, err
	}
	bf := &bytes.Buffer{}
	if err = t.template.Execute(bf, data); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (t *Template) trim(text string) (string, error) {
	var err error
	if text, err = t.convertDot(text); err != nil {
		return "", err
	}
	if text, err = t.trimMark(text); err != nil {
		return "", err
	}
	return text, nil
}

func (t *Template) convertDot(text string) (string, error) {
	reg, err := regexp.Compile(t.left + ".*?" + t.right)
	if  err != nil {
		return "", err
	}
	return reg.ReplaceAllStringFunc(text, func(s string) string {
		return strings.ReplaceAll(s, ".", "")
	}), nil
}

func (t *Template) trimMark(text string) (string, error) {
	reg, err := regexp.Compile("\"" + t.left + ".*?" + t.right + "\"")
	if  err != nil {
		return "", err
	}
	return reg.ReplaceAllStringFunc(text, func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\"", ""), "-", "")
	}), nil
}
