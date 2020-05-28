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

var ctx context.Context
var client *tts.Client
var encoding ttsapi.AudioEncoding
var voice ttsapi.SsmlVoiceGender
var language string

//Authenticate creates an internal client, authenticated with the given JSON credentials file or by environment variable (GOOGLE_APPLICATION_CREDENTIALS).
func Authenticate(credentialsFilepath string) {
	ctx = context.Background()
	c, err := tts.NewClient(ctx, option.WithCredentialsFile(credentialsFilepath))
	//Exit if client could not be created with supplied credentials.
	if err != nil {
		log.Fatalf("voice:Initialize() %v", err)
	}
	client = c
	SetSynthOptions("mp3", "N", "en-US")
}

//Synthesize synchronously converts text to audio and saves to file. May modify filename extension to match audio encoding.
func Synthesize(text, filename string) (outFilename string, err error) {
	timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	request := createRequest(text)
	response, err := client.SynthesizeSpeech(timeout, &request)
	if err != nil {
		log.Printf("voice:Synthesize() %v", err)
	} else {
		outFilename, err = saveAudioToFile(response.AudioContent, filename)
	}
	return
}

//SetSynthOptions sets options used for converting text to audio, such as encoding(mp3/wav/ogg), gender(M/F/N) and language(en-US).
func SetSynthOptions(audioEncoding, voiceGender, languageCode string) {
	setEncoding(audioEncoding)
	setVoice(voiceGender)
	language = languageCode
}

//Stop closes the client connection to the API service.
func Stop() (err error) {
	err = client.Close()
	if err != nil {
		log.Printf("voice:Stop() %v", err)
	}
	return
}

func createRequest(text string) (req ttsapi.SynthesizeSpeechRequest) {
	req = ttsapi.SynthesizeSpeechRequest{
		Input: &ttsapi.SynthesisInput{
			InputSource: &ttsapi.SynthesisInput_Text{Text: text},
		},
		Voice: &ttsapi.VoiceSelectionParams{
			LanguageCode: language,
			SsmlGender:   voice,
		},
		AudioConfig: &ttsapi.AudioConfig{
			AudioEncoding: encoding,
		},
	}
	return
}

func setEncoding(audioEncoding string) {
	enc := strings.ToUpper(audioEncoding)
	switch enc {
	case "OGG":
		encoding = ttsapi.AudioEncoding_OGG_OPUS
	case "WAV":
		encoding = ttsapi.AudioEncoding_LINEAR16
	default:
		encoding = ttsapi.AudioEncoding_MP3
	}
}

func setVoice(gender string) {
	vce := strings.ToUpper(gender)[0]
	switch string(vce) {
	case "M":
		voice = ttsapi.SsmlVoiceGender_MALE
	case "F":
		voice = ttsapi.SsmlVoiceGender_FEMALE
	default:
		voice = ttsapi.SsmlVoiceGender_NEUTRAL
	}
}

func saveAudioToFile(audio []byte, filename string) (file string, err error) {
	file = formatAudioFilename(filename)
	err = ioutil.WriteFile(file, audio, 0644)
	if err != nil {
		log.Printf("voice:saveAudioToFile() %v", err)
	}
	return
}

func formatAudioFilename(filename string) (formattedName string) {
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
			formattedName = filename[:extIndex] + expectedExt
		} else {
			formattedName = filename
		}
	} else {
		formattedName = filename + expectedExt
	}
	return
}
