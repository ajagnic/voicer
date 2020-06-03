package web

import (
	"net/http"
	"strconv"
)

func (c *webClient) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, c.htmlfile)

	} else {
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
	}
}

func (c *webClient) postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		filename := strconv.Itoa(c.UID)
		c.UID++
		r.ParseForm()
		input := r.Form.Get("input")
		c.AudioFile, _ = c.Synthesize(input, filename)
		generateHTML(*c)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		http.Error(w, "Invalid method.", http.StatusMethodNotAllowed)
	}
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/media"):]
	http.ServeFile(w, r, "."+path)
}
