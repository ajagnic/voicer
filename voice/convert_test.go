package voice

import (
	"flag"
	"os"
	"testing"

	ttsapi "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var key = flag.String("key", "", "Filepath of GCP Service-Account key")

func TestClient(t *testing.T) {
	Authenticate(*key)
	if client == nil {
		t.Errorf("Authenticate failed to create client. client=%v", client)
	} else {
		defer client.Close()
	}

	t.Run("SetSynthOptions", testSetSynthOptions)
	t.Run("Synthesize", testSynthesize)
}

func testSetSynthOptions(t *testing.T) {
	enc := "wav"
	vce := "Female"
	lang := "en-GB"
	SetSynthOptions(enc, vce, lang)
	if language != lang || voice != ttsapi.SsmlVoiceGender_FEMALE || encoding != ttsapi.AudioEncoding_LINEAR16 {
		t.Errorf("SetSynthOptions value mismatch. %v, %v, %v", language, voice, encoding)
	}
	SetSynthOptions("mp3", "N", "en-US")
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
