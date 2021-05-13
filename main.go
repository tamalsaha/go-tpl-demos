package main

import (
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"os"
	"text/template"
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

func main() {
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
