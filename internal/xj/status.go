package xj

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type state struct {
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	Descr  string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (s *state) Str() string {
	res, _ := jsoniter.MarshalIndent(s, "", "    ")
	return string(res)
}

func (s *state) Bts() []byte {
	res, _ := jsoniter.MarshalIndent(s, "", "    ")
	return res
}

func Info(v ...any) *state {
	return &state{Status: "info", Descr: fmt.Sprint(v...)}
}

func Infof(f string, v ...any) *state {
	return &state{Status: "info", Descr: fmt.Sprintf(f, v...)}
}

func Err(v ...any) *state {
	return &state{Status: "error", Descr: fmt.Sprint(v...)}
}

func Errf(f string, v ...any) *state {
	return &state{Status: "error", Descr: fmt.Sprintf(f, v...)}
}

func Succes(v ...any) *state {
	return &state{Status: "succes", Descr: fmt.Sprint(v...)}
}

func Succesf(f string, v ...any) *state {
	return &state{Status: "succes", Descr: fmt.Sprintf(f, v...)}
}

func Warn(v ...any) *state {
	return &state{Status: "warning", Descr: fmt.Sprint(v...)}
}

func Warnf(f string, v ...any) *state {
	return &state{Status: "warning", Descr: fmt.Sprintf(f, v...)}
}
