package voice

import (
	"flag"
	"os"
	"testing"

	ttsapi "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var keyFile = flag.String("key", "", "Filepath of GCP credentials file.")

func TestClient(t *testing.T) {
	Authenticate(*keyFile)
	if client == nil {
		t.Errorf("Authenticate failed to create client. client=%v", client)
	} else {
		defer client.Close()
	}

	t.Run("SetSynthOptions", testSetSynthOptions)
	t.Run("Synthesize", testSynthesize)
}

func testSetSynthOptions(t *testing.T) {
	lang := "en-GB"
	vce := "Female"
	enc := "wav"
	SetSynthOptions(lang, vce, enc)
	if language != lang || voice != ttsapi.SsmlVoiceGender_FEMALE || encoding != ttsapi.AudioEncoding_LINEAR16 {
		t.Errorf("SetSynthOptions value mismatch. %v, %v, %v", language, voice, encoding)
	}
	SetSynthOptions("en-US", "N", "mp3")
}

func testSynthesize(t *testing.T) {
	msg := "testing"
	filename := "output.wav"
	correctFilename := "output.mp3"
	outputFile, err := Synthesize(msg, filename)
	if err != nil {
		t.Errorf("Synthesize returned error. %v", err)
	}
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Synthesize file not found. %v", err)
	}
	if outputFile != correctFilename {
		t.Errorf("Synthesize unexpected filename. outputFile=%v", outputFile)
	}
	os.Remove(outputFile)
}