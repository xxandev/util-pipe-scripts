package main

import (
	"fmt"
	"util-pipe/internal/xj"
)

type UI64 uint64
type F64 float64
type Obj struct {
	TypeString string  `json:"type_string,omitempty" yaml:"type_string,omitempty" default:"text text text"`
	TypeInt    int     `json:"type_int,omitempty" yaml:"type_int,omitempty" default:"-3"`
	TypeUInt   UI64    `json:"type_uint,omitempty" yaml:"type_uint,omitempty" default:"33333333"`
	TypeF64    float64 `json:"type_f64,omitempty" yaml:"type_f64,omitempty" default:"30.5"`
	TypeBool   bool    `json:"type_bool,omitempty" yaml:"type_bool,omitempty" default:"true"`
}

type Config struct {
	TypeStruct Obj
	TypeStr    string             `default:"text text text"`
	TypeInt    int                `default:"-3"`
	TypeUInt   uint               `default:"33333333"`
	TypeUI64   UI64               `default:"5555"`
	TypeF64    F64                `default:"30.5"`
	TypeBool   bool               `default:"true"`
	TypeS1     []string           `default:"link1,link2"`
	TypeS2     []int              `default:"1,12,22,-37"`
	TypeS3     []int16            `default:"1,12,22,-37"`
	TypeS4     []float32          `default:"3.5,-30.5,22.0,-37"`
	TypeS5     []float64          `default:"3.5"`
	TypeM1     map[string]string  `default:"key1:value1,key2:value2,key3:value3"`
	TypeM2     map[string]int     `default:"key1:44,key2:55"`
	TypeM3     map[string]bool    `default:"key1:true,key2:true"`
	TypeM4     map[string]float64 `default:"key1:32.1,key2:77.25,key3:89.0"`
	TypeM5     map[string]float64 `default:"key1:32.1"`
	TypeAny1   any                `default:"-3"`
	TypeAny2   any                `default:"true"`
	TypeAny3   any                `default:"3.5,-30.5,22.0,-37"`
	TypeAny4   any                `default:"key1:32.1,key2:77.25,key3:89.0"`
	TypeAny5   any                `default:"key1:value1,key2:value2,key3:value3"`
}

func main() {
	var cfg Config
	var cfgINT int = 32
	fmt.Println(xj.SetDefaults(&cfg))
	fmt.Println(xj.SetDefaults(&cfgINT))
	fmt.Println("===================================")
	fmt.Printf("%#v\n", cfg)
	fmt.Println("===================================")
	fmt.Printf("%v\n", cfg)
}
