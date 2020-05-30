package voice

import (
	"context"
	"flag"
	"os"
	"testing"

	tts "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	ttsapi "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var key = flag.String("key", "", "Filepath of GCP Service-Account key")

func TestAuthenticate(t *testing.T) {
	client, err := Authenticate(*key)
	if err != nil {
		t.Errorf("Authenticate failed to create client. client=%v, err=%v", client, err)
	}
	err = client.Close()
	if err != nil {
		t.Errorf("Error closing the client. %v", err)
	}
}

func TestSynthesize(t *testing.T) {
	ttsClient, err := tts.NewClient(context.Background(), option.WithCredentialsFile(*key))
	if err != nil {
		t.Errorf("Error creating a client. %v", err)
	}
	client := &voiceClient{
		ttsClient,
		ttsapi.AudioEncoding_MP3,
		ttsapi.SsmlVoiceGender_NEUTRAL,
		"en-US",
	}
	filename := "output.wav"
	correctFilename := "output.mp3"
	outfile, err := client.Synthesize("testing", filename)
	if err != nil {
		t.Errorf("Synthesize returned error. %v", err)
	}
	if _, err = os.Stat(outfile); os.IsNotExist(err) {
		t.Errorf("Synthesize file not found. %v", err)
	}
	if outfile != correctFilename {
		t.Errorf("Synthesize unexpected filename. %v!=%v", outfile, correctFilename)
	}
	os.Remove(outfile)
	err = client.Close()
	if err != nil {
		t.Errorf("Error closing the client. %v", err)
	}
}

func TestSetSynthOptions(t *testing.T) {
	ttsClient, err := tts.NewClient(context.Background(), option.WithCredentialsFile(*key))
	if err != nil {
		t.Errorf("Error creating a client. %v", err)
	}
	client := &voiceClient{
		ttsClient,
		ttsapi.AudioEncoding_MP3,
		ttsapi.SsmlVoiceGender_NEUTRAL,
		"en-US",
	}
	client.SetSynthOptions("wav", "Female", "en-GB")
	if client.encoding != ttsapi.AudioEncoding_LINEAR16 || client.voice != ttsapi.SsmlVoiceGender_FEMALE || client.language != "en-GB" {
		t.Errorf("SetSynthOptions value mismatch. %v, %v, %v", client.encoding, client.voice, client.language)
	}
	err = client.Close()
	if err != nil {
		t.Errorf("Error closing the client. %v", err)
	}
}
