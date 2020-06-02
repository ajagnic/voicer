package web

import (
	"net/http"
)

func (c *webClient) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, c.htmlfile)

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		input := r.Form.Get("input")
		c.AudioFile, _ = c.Synthesize(input, c.AudioFile)
		http.ServeFile(w, r, c.htmlfile)

	} else {
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
	}
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/media"):]
	http.ServeFile(w, r, "."+path)
}
