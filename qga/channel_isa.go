package qga

import (
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

type IsaChannel struct {
	path string

	f *os.File

	req chan *Request
	res chan *Response

	m sync.Mutex
}

func NewIsaChannel(path string) (*IsaChannel, error) {
	return &IsaChannel{path: path}, nil
}

func (ch *IsaChannel) Open() error {
	var f *os.File
	var err error

	openMode := os.FileMode(os.ModeCharDevice | 0600)
	openFlags := unix.O_RDWR | ^unix.O_NONBLOCK | unix.O_CLOEXEC | unix.O_NDELAY

	if f, err = os.OpenFile(ch.path, openFlags, openMode); err != nil {
		return err
	}
	ch.f = f
	ch.req = make(chan *Request)
	ch.res = make(chan *Response, 1)

	return err
}

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

func (ch *IsaChannel) Close() error {
	if err := ch.f.Close(); err != nil {
		return err
	}
	close(ch.req)
	close(ch.res)
	return nil
}

func (ch *IsaChannel) Read(b []byte) (int, error) {
	return 0, nil
}

func (ch *IsaChannel) Write(b []byte) (int, error) {
	return 0, nil
}
