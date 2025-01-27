package yabin

import "errors"

var (
	// ErrAPIURLNotSet is returned when a call to the API is
	// attempted before a target instance URL has been set.
	ErrAPIURLNotSet = errors.New("YABin API URL not set")

	// ErrAPIClientNotSet is returned when the library is called
	// without a valid HTTP client.
	ErrAPIClientNotSet = errors.New("HTTP client not set")
	// ErrInvalidKey is returned when trying to search for a paste
	// with an invalid paste key / ID.
	ErrInvalidKey = errors.New("invalid paste key")

	// ErrInvalidContent is returned when the given paste content
	// is not valid.
	ErrInvalidContent = errors.New("invalid paste content")

	// ErrEncNotImplemented is returned when an operation demanding
	// encryption is requested.
	ErrEncNotImplemented = errors.New("encryption is not yet implemented in this library")
)
