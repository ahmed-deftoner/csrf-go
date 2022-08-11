package templates

import "text/template"

var templates = template.Must(template.ParseFiles("./server/templates/templateFiles/login.tmpl", "./server/templates/templateFiles/register.tmpl", "./server/templates/templateFiles/restricted.tmpl"))
