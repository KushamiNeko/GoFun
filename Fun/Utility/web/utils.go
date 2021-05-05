package web

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/KushamiNeko/GoFun/Utility/pretty"
)

func WriteTemplate(w http.ResponseWriter, temp *template.Template, name string, data interface{}) {
	buffer := bytes.Buffer{}

	err := temp.ExecuteTemplate(
		&buffer,
		name,
		data,
	)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperPink300, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(CleanAll(buffer.Bytes()))
	if err != nil {
		pretty.ColorPrintln(pretty.PaperPink300, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
