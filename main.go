package main

import (
	"fmt"
	"os"
	"text/template"
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

	tpl.Execute(os.Stdout, &na)
}