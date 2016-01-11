// +build linux

package channel

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// Poll channel for new messages
func (ch *VirtioChannel) Poll( /*block bool*/ ) error {
	var err error

	fd := int(ch.f.Fd())

	if err = unix.SetNonblock(fd, true); err != nil {
		return err
	}

	ch.pfd, err = unix.EpollCreate(1)
	if err != nil {
		return err
	}

	ctlEvent := unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLHUP | unix.EPOLLRDHUP, Fd: int32(fd)}
	if err = unix.EpollCtl(ch.pfd, unix.EPOLL_CTL_ADD, fd, &ctlEvent); err != nil {
		return err
	}
	events := make([]unix.EpollEvent, 32)

	for {
		nevents, err := unix.EpollWait(ch.pfd, events, 1000*60*5)
		switch err {
		case nil:
			if nevents == 0 {
				return fmt.Errorf("i/o timeout")
			}
			for ev := 0; ev < nevents; ev++ {
				if events[ev].Events == 0 {
					continue
				}
				if events[ev].Events&unix.EPOLLIN != 0 {
					if err = Serve(ch); err != nil {
						return err
					}
				}
			}
		case unix.EINTR:
			//			return
			continue
		default:
			return err
		}
	}

}
