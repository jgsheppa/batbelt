package batbelt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
)

type Batbelt struct {
	// Reader is the underlying reader.
	Reader         ReadAutoCloser
	stdout, stderr io.Writer
	httpClient     *http.Client

	// because pipe stages are concurrent, protect 'err'
	mu  *sync.Mutex
	err error
}

// ReadAutoCloser wraps an [io.ReadCloser] so that it will be automatically
// closed once it has been fully read.
type ReadAutoCloser struct {
	r io.ReadCloser
}

func NewBatbelt() *Batbelt {
	return &Batbelt{
		Reader:     ReadAutoCloser{},
		mu:         new(sync.Mutex),
		stdout:     os.Stdout,
		httpClient: http.DefaultClient,
	}
}

// SetError sets the error err on the batbelt.
func (b *Batbelt) SetError(err error) {
	if b.mu == nil { // uninitialised pipe
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.err = err
}

// WithError sets the error err on the batbelt.
func (b *Batbelt) WithError(err error) *Batbelt {
	b.SetError(err)
	return b
}

// Error returns any error present on the belt, or nil otherwise.
func (b *Batbelt) Error() error {
	if b.mu == nil { // uninitialised belt
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.err
}

// RemoveFile checks if a file exists and removes it
// if a file is found
func (b *Batbelt) RemoveFile(filepath string) *Batbelt {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return b
	}
	if err := os.Remove(filepath); err != nil {
		b.SetError(err)
	}

	return b
}

// CreateJSONFile is used to generate a JSON file
// in the file system given a filepath and data of
// any format
func (b *Batbelt) CreateJSONFile(structure any, filename string) *Batbelt {
	jsonBytes, err := json.Marshal(structure)
	if err != nil {
		b.SetError(err)
	}

	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		b.SetError(err)
	}
	return nil
}
