// Package yabin provides a library for communicating with YABin pastebin servers.
// Note: This package does not yet support encryption.
package yabin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	enry "github.com/go-enry/go-enry/v2"
	enryData "github.com/go-enry/go-enry/v2/data"
)

var apiURL string
var apiToken string

// POST to write, GET to read.
const apiRequestURI string = "/api/paste"

func buildAPIURL() string {
	return fmt.Sprintf("%s%s", apiURL, apiRequestURI)
}

// httpClient is the default HTTP Client used to make requests. You can override this with SetClient().
var httpClient = &http.Client{}

func isReady() error {
	if apiURL == "" {
		return ErrAPIURLNotSet
	}
	if httpClient == nil {
		return ErrAPIClientNotSet
	}
	return nil
}

var langCandidates []string = make([]string, 0, len(enryData.LanguagesLogProbabilities))

func populateLanguageCandidates() {
	for lang := range enryData.LanguagesLogProbabilities {
		langCandidates = append(langCandidates, lang)
	}
}

func guessContentLanguage(content string) (lang string) {
	if len(langCandidates) == 0 {
		populateLanguageCandidates()
	}
	lang, _ = enry.GetLanguageByClassifier([]byte(content), langCandidates)
	return lang
}

func paste(content, language, desiredKey string, expiry time.Duration, encPass *string, burnAfterRead bool) (key string, err error) {
	err = isReady()
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", ErrInvalidContent
	}

	// TODO Implement encryption
	// https://github.com/Yureien/YABin/blob/main/src/lib/crypto.ts
	if encPass != nil {
		return "", ErrEncNotImplemented
	}

	// Override language to plaintext if it's not given.
	if language == "" {
		language = "plaintext"
	}

	var expiresAfter *int
	if expiry.Seconds() > 0 {
		// FIXME possible bad scoping here / expiresSeconds may be dereferenced
		expiresSeconds := int(expiry.Seconds())
		expiresAfter = &expiresSeconds
	}

	// Build request
	reqData := WriteRequest{
		Content: content,
		// FIXME This feels ewwy. Maybe move to a subStruct?
		Config: struct {
			Language      *string `json:"language,omitempty"`
			Encrypted     bool    `json:"encrypted"`
			ExpiresAfter  *int    `json:"expiresAfter,omitempty"`
			BurnAfterRead bool    `json:"burnAfterRead"`
			CustomPath    string  `json:"customPath,omitempty"`
		}{
			Language:      &language,
			Encrypted:     false,
			ExpiresAfter:  expiresAfter,
			BurnAfterRead: burnAfterRead,
			CustomPath:    desiredKey,
		},
		PasswordProtected: false,
		InitVector:        "",
	}

	requestJson, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", buildAPIURL(), bytes.NewBuffer(requestJson))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	// No need to use a cookieJar, we're just doing one thing.
	if apiToken != "" {
		req.Header.Add("Cookie", "token="+apiToken)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := new(APIError)
		apiErr.StatusCode = resp.StatusCode
		if err = json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			apiErr.Message = err.Error()
			return "", apiErr
		}

		return "", apiErr
	}

	result := new(ReadResponse)
	err = json.NewDecoder(resp.Body).Decode(result)

	return result.Data.Key, err
}
