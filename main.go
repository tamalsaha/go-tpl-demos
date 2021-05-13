package main

import (
	"fmt"
	"os"
	"text/template"
)

var c1 = `
{{ define "t2" }}
{{ printf "%v" .A }}
{{ end }}
`

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
		},
	}
	tpl := template.Must(template.New("").Parse(`
{{ define "t2" }}
{{ printf "%v" .A }}
{{ end }}
{{ range $svc := .O }}
	{{ template "t2" $svc }}
{{ end }}
`))
	for _, tt := range tpl.Templates() {
		fmt.Println(tt.Name())
	}

	err := tpl.Execute(os.Stdout, &na)
	if err != nil {
		panic(err)
	}
}
