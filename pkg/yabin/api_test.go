package yabin

import (
	"net/http"
	"testing"
)

// TODO integration testing

func TestSetURL(t *testing.T) {
	// Test that the trailing forwardslash is stripped properly.
	SetURL("https://test.com/")

	if apiURL != "https://test.com" {
		t.Errorf("apiURL = %v, want %v", apiURL, "https://test.com")
	}
}

func TestSetAPIToken(t *testing.T) {
	SetAPIToken("token")
	if apiToken != "token" {
		t.Errorf("apiToken = %v, want %v", apiToken, "token")
	}
}

func TestSetClient(t *testing.T) {
	client := http.Client{}
	SetClient(&client)
	if httpClient != &client {
		t.Errorf("httpClient = %v, want %v", httpClient, client)
	}
}
