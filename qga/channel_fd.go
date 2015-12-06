package qga

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	fdChannelMode  = os.FileMode(0600)
	fdChannelFlags = unix.O_RDWR | ^unix.O_NONBLOCK | unix.O_CLOEXEC | unix.O_NDELAY
)

// FdChannel struct
type FdChannel struct {
	path string

	f *os.File
	r *bufio.Reader

	req chan *Request
	res chan *Response

	m sync.Mutex
}

// NewFdChannel create new fd channel
func NewFdChannel(path string) (*FdChannel, error) {
	return &FdChannel{path: path}, nil
}

// Open opens connection
func (ch *FdChannel) Open() error {
	var f *os.File
	var err error

	fd, err := strconv.Atoi(ch.path)
	if err != nil || fd < 0 {
		return fmt.Errorf("invalid fd")
	}

	f = os.NewFile(uintptr(fd), ch.path)

	ch.f = f
	ch.r = bufio.NewReaderSize(f, QGA_MAX_MESSAGE_LEN)
	ch.req = make(chan *Request)
	ch.res = make(chan *Response, 1)

	return err
}

// Reset reinit channel
func (ch *FdChannel) Reset() error {
	var err error
	ch.m.Lock()
	defer ch.m.Unlock()
	err = ch.Close()
	if err != nil {
		return err
	}
	return ch.Open()
}

// Close closes channel
func (ch *FdChannel) Close() error {
	if err := ch.f.Close(); err != nil {
		return err
	}
	close(ch.req)
	close(ch.res)
	return nil
}

// Read read from channel
func (ch *FdChannel) Read(buffer []byte) (int, error) {
	if ch.f == nil {
		return 0, fmt.Errorf("try to read on closed channel")
	}
	return ch.f.Read(buffer)
}

// Write wtire to channel
func (ch *FdChannel) Write(buffer []byte) (int, error) {
	if ch.f == nil {
		return 0, fmt.Errorf("try to write on closed channel")
	}
	return ch.f.Write(buffer)
}
