package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gomodules.xyz/encoding/json"
	"gomodules.xyz/encoding/yaml"
	"k8s.io/klog/v2"
)

var (
	c1 = `
{{ define "t1" }}
{{ printf "c1-t1-%v" .A }}
{{ end }}
{{ define "t2" }}
{{ printf "c1-t2-%v" .A }}
{{ end }}
`
	c2 = `
{{ define "t1" }}
{{ printf "c2-t1-%v" .A }}
{{ end }}
{{ define "t2" }}
{{ printf "c2-t2-%v" .A }}
{{ end }}
`
	content = `
{{ range $svc := .O }}
	{{ template "t2" $svc }}
{{ end }}
`
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	data, err := ioutil.ReadFile("pod.yaml")
	if err != nil {
		klog.Fatalln(err)
	}

	var obj map[string]interface{}
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		klog.Fatalln(err)
	}

	funcs := sprig.TxtFuncMap()
	funcs["custom_int"] = tplCustomIntFn
	funcs["custom_str"] = tplCustomStrFn
	funcs["custom_obj"] = tplCustomObjFn
	funcs["custom_struct"] = tplCustomStructFn

	// txt := `{{ toRawJson . | custom_obj | toRawJson }}`
	// txt := `{{ custom_obj . | toRawJson }}`
	// txt := `{{ custom_int . }}`
	// txt := `{{ custom_str . }}`
	txt := `{{ custom_struct . | toRawJson }}`
	tpl, err := template.New("").Funcs(funcs).Parse(txt)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(os.Stdout, obj)
	if err != nil {
		panic(err)
	}
}

func tplCustomIntFn(data interface{}) (int64, error) {
	return 1, nil
}

func tplCustomStrFn(data interface{}) (string, error) {
	return "abc", nil
}

func tplCustomObjFn(data interface{}) (interface{}, error) {
	return toObject(data)
}

func tplCustomStructFn(data interface{}) (Person, error) {
	return Person{
		Name: "John",
		Age:  30,
	}, nil
}

func toObject(data interface{}) (map[string]interface{}, error) {
	var obj map[string]interface{}
	if v, ok := data.(map[string]interface{}); ok {
		obj = v
	} else if str, ok := data.(string); ok {
		err := json.Unmarshal([]byte(str), &obj)
		if err != nil {
			return nil,  err
		}
	} else {
		return nil, fmt.Errorf("unknown obj type %v", reflect.TypeOf(data).String())
	}
	return obj,  nil
}


func main__() {
	type Inner struct {
		A string
	}
	type Outer struct {
		Inner
	}

	type NA struct {
		O []Outer
	}

	na := NA{
		O: []Outer{
			{
				Inner: Inner{A: "123"},
			},
			{
				Inner: Inner{A: "456"},
			},
		},
	}

	tpl_c1, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(c1)
	if err != nil {
		panic(err)
	}
	tpl_c2, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(c2)
	if err != nil {
		panic(err)
	}

	m_c1 := template.New("")
	for _, tt := range tpl_c1.Templates() {
		if tt.Name() != "" {
			_, err = m_c1.AddParseTree("c1_"+tt.Name(), tt.Tree) // m_c1
			if err != nil {
				panic(err)
			}
			_, err = m_c1.AddParseTree(tt.Name(), tt.Tree) // m_c1
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(tt.Name())
		}
	}
	_, err = m_c1.Parse(content)
	if err != nil {
		panic(err)
	}

	m_c2 := template.New("")
	for _, tt := range tpl_c2.Templates() {
		if tt.Name() != "" {
			_, err = m_c2.AddParseTree("c2_"+tt.Name(), tt.Tree) // m_c2
			if err != nil {
				panic(err)
			}
			_, err = m_c2.AddParseTree(tt.Name(), tt.Tree) // m_c2
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(tt.Name())
		}
	}
	_, err = m_c2.Parse(content)
	if err != nil {
		panic(err)
	}

	err = m_c1.Execute(os.Stdout, &na)
	if err != nil {
		panic(err)
	}
	fmt.Println("=====================================")
	err = m_c2.Execute(os.Stdout, &na)
	if err != nil {
		panic(err)
	}
}
