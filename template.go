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
		for key, value := range m {
			if _, ok = value.(string); ok {
				//m[key] = strconv.QuoteToASCII(v)
				continue
			}
			var vBytes []byte
			if vBytes, err = json.Marshal(value); err != nil {
				return nil, err
			}
			m[key] = string(vBytes)
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
	if text, err = t.trimSpace(text); err != nil {
		return "", err
	}
	if text, err = t.trimMark(text); err != nil {
		return "", err
	}
	if text, err = t.convertDot(text); err != nil {
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
		tmp := []byte(s)
		dot := []byte(".")
		tmp = append(tmp[:2], append(dot, tmp[2:]...)...)
		return string(tmp)
	}), nil
}

func (t *Template) trimMark(text string) (string, error) {
	reg, err := regexp.Compile("\"" + t.left + "-.*?-" + t.right + "\"")
	if  err != nil {
		return "", err
	}
	return reg.ReplaceAllStringFunc(text, func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\"", ""), "-", "")
	}), nil
}

func (t *Template) trimSpace(text string) (string, error) {
	reg, err := regexp.Compile(t.left + ".*?" + t.right)
	if  err != nil {
		return "", err
	}
	return reg.ReplaceAllStringFunc(text, func(s string) string {
		return strings.ReplaceAll(s, " ", "")
	}), nil
}
