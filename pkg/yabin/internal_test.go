package yabin

import (
	"errors"
	"fmt"
	"testing"
)

func testResetState() {
	apiURL = ""
	apiToken = ""
}

func TestBuildAPIURL(t *testing.T) {
	t.Cleanup(testResetState)
	SetURL("http://example.com")
	if buildAPIURL() != fmt.Sprintf("http://example.com%s", apiRequestURI) {
		t.Errorf("BuildAPIURL returned unexpected URL %s", buildAPIURL())
	}
}

func TestReadiness(t *testing.T) {
	t.Cleanup(testResetState)
	if e := isReady(); !errors.Is(e, ErrAPIURLNotSet) {
		t.Errorf("isReady() returned unexpected error: %v", e)
	}

	SetURL("http://example.com")
	if e := isReady(); e != nil {
		t.Errorf("isReady() returned faults after initialisation: %v", e)
	}

	// intentionally break things
	SetClient(nil)
	if e := isReady(); !errors.Is(e, ErrAPIClientNotSet) {
		t.Errorf("isReady() returned unexpected error: %v", e)
	}
}
