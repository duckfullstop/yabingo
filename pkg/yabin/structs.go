package yabin

import "fmt"

type WriteRequest struct {
	// Paste content. If encrypted, must be encoded into a Base64 string.
	Content string `json:"content"`
	Config  struct {
		// Programming language of the paste. Defaults to plaintext.
		Language *string `json:"language,omitempty"`
		// Whether the paste is encrypted.
		Encrypted bool `json:"encrypted"`
		// Time in seconds until the paste expires.
		ExpiresAfter *int `json:"expiresAfter,omitempty"`
		// Whether the paste should be deleted after reading.
		BurnAfterRead bool `json:"burnAfterRead"`
		// A custom path for the paste.
		CustomPath string `json:"customPath,omitempty"`
	} `json:"config"`
	// Whether the paste should be password protected.
	PasswordProtected bool `json:"passwordProtected"`
	// Initialization vector for AES encryption. Max length: 64.
	InitVector string `json:"initVector,omitempty"`
}

type ReadResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Key               string `json:"key"`
		Content           string `json:"content"`
		Encrypted         bool   `json:"encrypted"`
		PasswordProtected bool   `json:"passwordProtected"`
		InitVector        string `json:"initVector"`
		Language          string `json:"language"`
		OwnerId           string `json:"ownerId"`
	} `json:"data"`
}

// APIError defines the YABin API Error response.
type APIError struct {
	Success    bool `json:"success"`
	StatusCode int
	Message    string `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: %d - %s", e.StatusCode, e.Message)
}
