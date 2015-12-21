// +build ignore

package channel

import (
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	isaChannelMode  = os.FileMode(os.ModeCharDevice | 0600)
	isaChannelFlags = unix.O_RDWR | ^unix.O_NONBLOCK | unix.O_CLOEXEC | unix.O_NDELAY
)

// IsaChannel struct
type IsaChannel struct {
	path string

	f *os.File

	req chan *Request
	res chan *Response

	m sync.Mutex
}

// NewIsaChannel creates new isa channel
func NewIsaChannel(path string) (*IsaChannel, error) {
	return &IsaChannel{path: path}, nil
}

// Open initialize channel
func (ch *IsaChannel) Open() error {
	var f *os.File
	var err error

	if f, err = os.OpenFile(ch.path, isaChannelFlags, isaChannelMode); err != nil {
		return err
	}
	ch.f = f
	ch.req = make(chan *Request)
	ch.res = make(chan *Response, 1)

	return err
}

// Reset resets channel
func (ch *IsaChannel) Reset() error {
	var err error
	ch.m.Lock()
	defer ch.m.Unlock()
	err = ch.Close()
	if err != nil {
		return err
	}
	return ch.Open()
}

// Close closing channel
func (ch *IsaChannel) Close() error {
	if err := ch.f.Close(); err != nil {
		return err
	}
	close(ch.req)
	close(ch.res)
	return nil
}

// Read read from channel
func (ch *IsaChannel) Read(b []byte) (int, error) {
	return 0, nil
}

// Write write to channel
func (ch *IsaChannel) Write(b []byte) (int, error) {
	return 0, nil
}
