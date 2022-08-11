package templates

import "text/template"

type Register struct {
	BalertUser bool
	AlertMsg   string
}

type Login struct {
	BalertUser bool
	AlertMsg   string
}

var templates = template.Must(template.ParseFiles("./server/templates/templatefiles/login.tmpl", "./server/templates/templatefiles/register.tmpl", "./server/templates/templatefiles/restricted.tmpl"))
