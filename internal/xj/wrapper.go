package xj

import (
	"errors"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type wrapper struct {
	sep  string
	Data sync.Map
}

func NewWrap() *wrapper {
	return &wrapper{sep: "."}
}

func (wrp *wrapper) SetSeparator(v string) {
	wrp.sep = v
}

func (wrp *wrapper) Unmarshal(data []byte) error {
	var object any
	if err := jsoniter.Unmarshal(data, &object); err != nil {
		return err
	}
	return wrp.absorb("", object)
}

func (wrp *wrapper) Marshal() ([]byte, error) {
	object, err := wrp.extract()
	if err != nil {
		return nil, err
	}
	return jsoniter.Marshal(object)
}

func (wrp *wrapper) MarshalIndent(prefix string, indent string) ([]byte, error) {
	object, err := wrp.extract()
	if err != nil {
		return nil, err
	}
	return jsoniter.MarshalIndent(object, prefix, indent)
}

func (wrp *wrapper) absorb(prefix string, object any) error {
	if obj, ok := object.(map[string]any); ok {
		for key, val := range obj {
			if sobj, ok := val.(map[string]any); ok {
				wrp.absorb(prefix+key+wrp.sep, sobj)
				continue
			}
			wrp.Data.Store(prefix+key, val)
		}
		return nil
	}
	return errors.New("unknown data object")
}

func (wrp *wrapper) extract() (any, error) {
	result := make(map[string]any)
	wrp.Data.Range(func(key, value any) bool {
		vkey, ok := key.(string)
		if ok {
			skeys := strings.Split(vkey, wrp.sep)
			cmap := result
			for i := 0; i < len(skeys)-1; i++ {
				if _, ok := cmap[skeys[i]]; !ok {
					cmap[skeys[i]] = make(map[string]any)
				}
				cmap = cmap[skeys[i]].(map[string]any)
			}
			cmap[skeys[len(skeys)-1]] = value
		}
		return ok
	})
	return result, nil
}
