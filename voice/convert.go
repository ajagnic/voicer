package voice

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"time"

	tts "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	ttsapi "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type VoiceClient struct {
	*tts.Client
	encoding ttsapi.AudioEncoding
	voice    ttsapi.SsmlVoiceGender
	language string
}

//Authenticate returns a client, authenticated with the given JSON credentials file or by environment variable (GOOGLE_APPLICATION_CREDENTIALS).
func Authenticate(filepath string) (*VoiceClient, error) {
	ttsClient, err := tts.NewClient(context.Background(), option.WithCredentialsFile(filepath))
	if err != nil {
		log.Printf("voice:Authenticate() %v", err)
	}
	client := &VoiceClient{
		ttsClient,
		ttsapi.AudioEncoding_MP3,
		ttsapi.SsmlVoiceGender_NEUTRAL,
		"en-US",
	}
	return client, err
}

//Synthesize synchronously converts text to audio and saves to file. May modify filename extension to match audio encoding.
func (c *VoiceClient) Synthesize(text, filename string) (outfile string, err error) {
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request := createRequest(c, text)
	response, err := c.SynthesizeSpeech(timeout, request)
	if err != nil {
		log.Printf("voice:Synthesize() %v", err)
	} else {
		outfile, err = saveAudioToFile(response.AudioContent, filename, c.encoding)
	}
	return
}

//SetSynthOptions sets options used for converting text to audio, such as encoding(mp3/wav/ogg), gender(M/F/N) and language(en-US).
func (c *VoiceClient) SetSynthOptions(audioEncoding, voiceGender, languageCode string) {
	c.setEncoding(audioEncoding)
	c.setVoice(voiceGender)
	c.language = languageCode
}

func (c *VoiceClient) setEncoding(audioEncoding string) {
	enc := strings.ToUpper(audioEncoding)
	switch enc {
	case "OGG":
		c.encoding = ttsapi.AudioEncoding_OGG_OPUS
	case "WAV":
		c.encoding = ttsapi.AudioEncoding_LINEAR16
	default:
		c.encoding = ttsapi.AudioEncoding_MP3
	}
}

func (c *VoiceClient) setVoice(gender string) {
	vce := strings.ToUpper(gender)[:1]
	switch vce {
	case "M":
		c.voice = ttsapi.SsmlVoiceGender_MALE
	case "F":
		c.voice = ttsapi.SsmlVoiceGender_FEMALE
	default:
		c.voice = ttsapi.SsmlVoiceGender_NEUTRAL
	}
}

func createRequest(c *VoiceClient, text string) *ttsapi.SynthesizeSpeechRequest {
	return &ttsapi.SynthesizeSpeechRequest{
		Input: &ttsapi.SynthesisInput{
			InputSource: &ttsapi.SynthesisInput_Text{Text: text},
		},
		Voice: &ttsapi.VoiceSelectionParams{
			LanguageCode: c.language,
			SsmlGender:   c.voice,
		},
		AudioConfig: &ttsapi.AudioConfig{
			AudioEncoding: c.encoding,
		},
	}
}

func saveAudioToFile(audio []byte, filename string, encoding ttsapi.AudioEncoding) (string, error) {
	file := formatAudioFilename(filename, encoding)
	err := ioutil.WriteFile(file, audio, 0644)
	if err != nil {
		log.Printf("voice:saveAudioToFile() %v", err)
	}
	return file, err
}

func formatAudioFilename(filename string, encoding ttsapi.AudioEncoding) string {
	fmtname := filename
	var expectedExt string
	switch encoding {
	case ttsapi.AudioEncoding_OGG_OPUS:
		expectedExt = ".ogg"
	case ttsapi.AudioEncoding_LINEAR16:
		expectedExt = ".wav"
	case ttsapi.AudioEncoding_MP3:
		expectedExt = ".mp3"
	}
	if extIndex := strings.Index(filename, "."); extIndex != -1 {
		fileExt := filename[extIndex:]
		if fileExt != expectedExt {
			fmtname = filename[:extIndex] + expectedExt
		}
	} else {
		fmtname = filename + expectedExt
	}
	return fmtname
}
