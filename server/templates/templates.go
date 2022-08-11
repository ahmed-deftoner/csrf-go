package templates

import (
	"log"
	"net/http"
	"text/template"
)

type Register struct {
	BalertUser bool
	AlertMsg   string
}

type Login struct {
	BalertUser bool
	AlertMsg   string
}

type Restricted struct {
	CsrfSecret    string
	SecretMessage string
}

var templates = template.Must(template.ParseFiles("./server/templates/templatefiles/login.tmpl", "./server/templates/templatefiles/register.tmpl", "./server/templates/templatefiles/restricted.tmpl"))

func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", p)
	if err != nil {
		log.Printf("Temlate error here: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
