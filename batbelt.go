package batbelt

import (
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
