package qga

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	vioChannelMode  = os.FileMode(os.ModeCharDevice | 0600)
	vioChannelFlags = unix.O_RDWR | ^unix.O_NONBLOCK | unix.O_CLOEXEC | unix.O_NDELAY
)

// VirtioChannel struct
type VirtioChannel struct {
	path string

	f *os.File
	r *bufio.Reader

	req chan *Request
	res chan *Response

	m sync.Mutex
}

// NewVirtioChannel creates new virtio channel
func NewVirtioChannel(path string) (*VirtioChannel, error) {
	return &VirtioChannel{path: path}, nil
}

// Open open channel
func (ch *VirtioChannel) Open() error {
	var f *os.File
	var err error

	if f, err = os.OpenFile(ch.path, vioChannelFlags, vioChannelMode); err != nil {
		return err
	}
	ch.f = f
	ch.r = bufio.NewReaderSize(f, QGA_MAX_MESSAGE_LEN)
	ch.req = make(chan *Request)
	ch.res = make(chan *Response, 1)

	return err
}

// Reset rests channel
func (ch *VirtioChannel) Reset() error {
	var err error
	ch.m.Lock()
	defer ch.m.Unlock()
	err = ch.Close()
	if err != nil {
		return err
	}
	return ch.Open()
}

// Close close channel
func (ch *VirtioChannel) Close() error {
	if err := ch.f.Close(); err != nil {
		return err
	}
	close(ch.req)
	close(ch.res)
	return nil
}

// Read read from channel
func (ch *VirtioChannel) Read(buffer []byte) (int, error) {
	if ch.f == nil {
		return 0, fmt.Errorf("try to read on closed channel")
	}
	return ch.f.Read(buffer)
}

// Write write to channel
func (ch *VirtioChannel) Write(buffer []byte) (int, error) {
	if ch.f == nil {
		return 0, fmt.Errorf("try to write on closed channel")
	}
	return ch.f.Write(buffer)
}

/*
func (ch *VirtioChannel) ServeIO() error {
	var wg sync.WaitGroup

	//	wg.Add(2)

	// channel buffer worker
	go func() {
		defer wg.Done()
		var err error

		for {
			select {
			case buffer := <-ch.r:
				if err = json.Unmarshal(&req); err == nil {
					ch.req <- &req
				} else {
					ch.res <- &Response{Error: &Error{Code: -1, Desc: fmt.Sprintf("invalid request %s", err.Error())}}
				}
			case buffer := <-ch.w:
				ch.f.Write(buffer)
			}
		}
	}()

	//reader
	go func() {
		defer wg.Done()

		buffer := make([]byte, 4*1024)
		var n int
		var err error
		if n, err = ch.Write(buffer); err == nil {
			w <- buffer[:n]
		} else {
			fmt.Printf(err.Error())
		}
	}()

	wg.Wait()

	return nil
}
*/
