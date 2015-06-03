package adifparser

import (
	"bytes"
	"io"
	"net/http"
)

type LOTWClient interface {
	// Implement reader
	Read([]byte) (int, error)
	Close() error
}

type lotwClientImpl struct {
	// Creds
	username string
	password string
	// HTTP state
	httpResponse *http.Response
	// Temporary read buffer
	buf []byte
	// TODO: state storage
}

// Create a new client
func NewLOTWClient(username, password string) *lotwClientImpl {
	client := &lotwClientImpl{}
	client.username = username
	client.password = password
	client.buf = make([]byte, 0, 1024)
	return client
}

// Read from socket
func (c *lotwClientImpl) Read(p []byte) (int, error) {
	if c.httpResponse == nil {
		if err := c.open(); err != nil {
			// TODO: better logging
			return 0, err
		}
	}
	usable := c.buf
	storage := make([]byte, 0, cap(usable)-len(usable))
	_, err := c.httpResponse.Body.Read(storage)
	if err != nil && err != io.EOF {
		copyWithoutLOTW(p, usable)
		return len(p), err
	}
	if cap(p) < len(usable) {
		usable = usable[:cap(p)]
	}
	if end := bytes.LastIndex(usable, []byte("<")); end > 0 {
		usable = usable[:end]
	}
	c.buf = c.buf[len(usable):]
	copyWithoutLOTW(p, usable)
	if err == io.EOF && len(c.buf) == 0 {
		return len(p), err
	}
	return len(p), nil
}

func copyWithoutLOTW(dst, src []byte) int {
	if end := bytes.Index(src, []byte("<APP_LoTW_EOF>")); end > -1 {
		src = src[:end]
	}
	return copy(dst, src)
}

func (c *lotwClientImpl) Close() error {
	if c.httpResponse != nil {
		err := c.httpResponse.Body.Close()
		c.httpResponse = nil
		return err
	}
	return nil
}

func (c *lotwClientImpl) open() error {
	return nil
}
