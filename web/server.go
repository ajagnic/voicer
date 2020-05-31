package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Serve(addr string) error {
	router := http.NewServeMux()
	server := createServer(addr, router)
	interrupt := shutdownListener()

	router.HandleFunc("/", indexHandler)

	go listen(server)
	<-interrupt
	err := shutdown(server)
	return err
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world."))
}

func listen(srv *http.Server) {
	log.Printf("Server listening at %v.", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("web:listen() %v", err)
	}
}

func shutdown(srv *http.Server) error {
	timeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := srv.Shutdown(timeout)
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
