package main

import (
	"flag"

	"github.com/ajagnic/voicer/voice"
)

var key = flag.String("key", "", "Filepath of GCP Service-Account key")

func main() {
	flag.Parse()
	voice.Authenticate(*key)
	defer voice.Stop()
}
