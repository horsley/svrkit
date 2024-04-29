package svrkit

import (
	"encoding/json"
	"os"
)

// JSONStore 通用的简单存储基于json格式
type JSONStore[T interface{}] struct {
	Filename string
	Data     T
}

func NewJSONStore[T interface{}]() *JSONStore[T] {
	return &JSONStore[T]{}
}

func (j *JSONStore[T]) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	j.Filename = file
	return d.Decode(&j.Data)
}

func (j *JSONStore[T]) Save() error {
	f, err := os.OpenFile(j.Filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	return e.Encode(j.Data)
}
