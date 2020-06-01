package web

import (
	"net/http"
	"text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index := indexPage()
	tmpl := template.Must(template.New("layout").Parse(htmlLayout()))
	if r.Method == http.MethodGet {
		tmpl.Execute(w, index)
	} else {
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
	}
}

func (s *webClient) postHandler(w http.ResponseWriter, r *http.Request) {
	index := indexPage()
	tmpl := template.Must(template.New("layout").Parse(htmlLayout()))
	if r.Method == http.MethodPost {
		r.ParseForm()
		input := r.Form.Get("input")
		s.Synthesize(input, s.filename)
		tmpl.Execute(w, index)
	} else {
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
	}
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/media"):]
	http.ServeFile(w, r, "."+path)
}
