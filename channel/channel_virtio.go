// +build linux windows

package channel

import (
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	vioChannelMode  = os.FileMode(os.ModeExclusive | os.ModeCharDevice | 0600)
	vioChannelFlags = os.O_RDWR | unix.O_NONBLOCK | unix.O_ASYNC
)

// VirtioChannel struct
type VirtioChannel struct {
	path string
	f    *os.File
	m    sync.Mutex
	pfd  int
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

	return err
}

// Reset rests channel
func (ch *VirtioChannel) Reset() error {
	ch.m.Lock()
	defer ch.m.Unlock()

	err := ch.Close()
	if err != nil {
		return err
	}

	return ch.Open()
}

// Read read from channel
func (ch *VirtioChannel) Read(b []byte) (int, error) {
	ch.m.Lock()
	defer ch.m.Unlock()
	return ch.f.Read(b)
}

func (ch *VirtioChannel) Write(b []byte) (int, error) {
	ch.m.Lock()
	defer ch.m.Unlock()
	return ch.f.Write(b)
}

// Close close channel
func (ch *VirtioChannel) Close() error {
	ch.m.Lock()
	defer ch.m.Unlock()
	err := unix.Close(ch.pfd)
	if err != nil {
		return err
	}
	return ch.f.Close()
}
