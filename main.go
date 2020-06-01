package main

import (
	"flag"
	"log"

	"github.com/ajagnic/voicer/web"
)

var key = flag.String("key", "", "Filepath of GCP Service-Account key")

func main() {
	flag.Parse()
	client, err := web.AuthenticateServer(*key)
	if err != nil {
		log.Fatalf("Could not authenticate server. %v", err)
	}
	client.Serve(":8080")
}
