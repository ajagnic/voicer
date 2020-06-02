package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/ajagnic/voicer/voice"
)

type webClient struct {
	*voice.VoiceClient
	AudioFile string
	htmlfile  string
}

func AuthenticateServer(key string) (*webClient, error) {
	voiceClient, err := voice.Authenticate(key)
	if err != nil {
		log.Printf("web:AuthenticateServer() %v", err)
	}
	client := &webClient{
		voiceClient,
		"output",
		"index.html",
	}
	return client, err
}

func (c *webClient) Serve(addr string) error {
	router := http.NewServeMux()
	server := createServer(addr, router)
	interrupt := shutdownListener()
	err := generateHTML(*c)

	router.HandleFunc("/", c.indexHandler)
	router.HandleFunc("/media/", mediaHandler)

	go listen(server)
	<-interrupt
	err = shutdown(server, c)
	return err
}

func listen(srv *http.Server) {
	log.Printf("Server listening at %v.", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("web:listen() %v", err)
	}
}

func shutdown(srv *http.Server, client *webClient) error {
	timeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := srv.Shutdown(timeout)
	os.Remove(client.AudioFile)
	os.Remove(client.htmlfile)
	client.Close()
	log.Printf("Server shutdown.")
	return err
}

func shutdownListener() (sigint chan os.Signal) {
	sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	return
}

func createServer(addr string, router *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func generateHTML(client webClient) error {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Voicer</title>
	</head>
	<body>
		<audio autoplay controls>
			<source src="media/{{.AudioFile}}.mp3" type="audio/mpeg">
			<source src="media/{{.AudioFile}}.ogg" type="audio/ogg">
			<source src="media/{{.AudioFile}}.wav" type="audio/wav">
		</audio>
		<form action="/post" method="post">
			<textarea name="input" id="inputArea" cols="30" rows="10"></textarea>
			<button type="submit">Convert</button>
		</form>
	</body>
	</html>
	`
	file, err := os.OpenFile(client.htmlfile, os.O_CREATE|os.O_WRONLY, 0644)
	tmpl := template.Must(template.New("interface").Parse(html))
	err = tmpl.Execute(file, client)
	return err
}
