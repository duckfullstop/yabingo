package yabin

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// SetURL sets the URL of your target YABin server.
// The URL MUST be fully qualified; e.g. https://yabin.yourho.st
func SetURL(url string) {
	apiURL = strings.TrimSuffix(url, "/")
}

// SetAPIToken sets the API key to use when communicating with the target YABin server.
func SetAPIToken(key string) {
	apiToken = key
}

// SetClient overrides the http.Client used to make API requests.
func SetClient(client *http.Client) {
	httpClient = client
}

// GetPaste attempts to fetch a paste using the given paste key / ID.
// If the paste is encrypted, that's up to you.
func GetPaste(key string) (result *ReadResponse, err error) {
	err = isReady()
	if err != nil {
		return nil, err
	}

	if key == "" {
		return nil, ErrInvalidKey
	}

	req, err := http.NewRequest("GET", buildAPIURL(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.Query().Set("key", key)

	// not sure if this is necessary? doesn't hurt
	req.URL.RawQuery = req.URL.Query().Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := new(APIError)
		if err = json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return nil, err
		}

		return nil, apiErr
	}

	result = new(ReadResponse)
	err = json.NewDecoder(resp.Body).Decode(result)

	return result, err

}

// Paste sends the given content to the YABin server
// with the language automatically determined using
// heuristics.
func Paste(content string, burnAfterRead bool) (key string, err error) {
	return PasteWithLanguageExpiry(content, guessContentLanguage(content), 0, burnAfterRead)
}

// PastePlaintext sends the given content to the YABin server
// without determining its language.
func PastePlaintext(content string, burnAfterRead bool) (key string, err error) {
	return PasteWithLanguage(content, "plaintext", burnAfterRead)
}

// PasteWithKey sends the given content to the YABin server
// with the given requested key.
func PasteWithKey(content, requestedKey string, burnAfterRead bool) (key string, err error) {
	return paste(content, guessContentLanguage(content), requestedKey, 0, nil, burnAfterRead)
}

// PasteWithKeyLanguage sends the given content to the YABin server
// with the given requested key & language.
func PasteWithKeyLanguage(content, requestedKey, language string, burnAfterRead bool) (key string, err error) {
	return paste(content, language, requestedKey, 0, nil, burnAfterRead)
}

// PasteWithExpiry sends the given content to the YABin server
// with a given expiry duration.
func PasteWithExpiry(content string, expiry time.Duration, burnAfterRead bool) (key string, err error) {
	return PasteWithLanguageExpiry(content, guessContentLanguage(content), expiry, burnAfterRead)
}

// PasteWithLanguage sends the given content to the YABin server
// with the language set to the given argument.
func PasteWithLanguage(content string, language string, burnAfterRead bool) (key string, err error) {
	return PasteWithLanguageExpiry(content, language, 0, burnAfterRead)
}

// PasteWithLanguageExpiry sends the given content to the YABin server
// with the language set to the given argument
// and the expiry set to the given duration.
func PasteWithLanguageExpiry(content string, language string, expiry time.Duration, burnAfterRead bool) (key string, err error) {
	return paste(content, language, "", expiry, nil, burnAfterRead)
}
