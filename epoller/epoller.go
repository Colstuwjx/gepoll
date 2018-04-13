package epoller

/*
   Thanks https://github.com/wolfeidau/epoller.
*/

import (
	"os"
	"syscall"
)

const (
	// MaxEpollEvents max events to queue
	MaxEpollEvents = 8

	// MaxReadSize maximum read size
	MaxReadSize = 1024
)

type EventHandler func(slice []byte, n int)
type Epoller struct {
	fd      int
	epfd    int
	handler EventHandler
}

func NewEpoller(handler EventHandler) *Epoller {
	return &Epoller{
		handler: handler,
	}
}

func (ep *Epoller) Open(device string) error {
	fd, err := syscall.Open(device, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	ep.fd = fd
	return nil
}

func (ep *Epoller) Close() error {
	syscall.Close(ep.fd)
	return nil
}

func (ep *Epoller) notify(eventFd int) {
	var buf [MaxReadSize]byte
	n, _ := syscall.Read(eventFd, buf[:])
	if n > 0 {
		ep.handler(buf[0:n], n)
	}
}

func (ep *Epoller) Dispatch() error {
	var event syscall.EpollEvent
	var events [MaxEpollEvents]syscall.EpollEvent

	if err := syscall.SetNonblock(ep.fd, true); err != nil {
		return err
	}

	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epfd)

	event.Events = syscall.EPOLLIN
	event.Fd = int32(ep.fd)
	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, ep.fd, &event); err != nil {
		return err
	}

	var nevents int
	for {

		nevents, err = syscall.EpollWait(epfd, events[:], -1)
		if err != nil {
			return err
		}

		for ev := 0; ev < nevents; ev++ {
			// dispatch this to avoid delays in processing
			ep.notify(int(events[ev].Fd))
		}

	}
	return nil
}
