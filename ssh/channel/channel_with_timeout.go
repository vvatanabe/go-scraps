package channel

import (
	"errors"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrTimeoutRead = errors.New("timeout read")

var ErrTimeoutWrite = errors.New("timeout write")

var errTimeout = errors.New("timeout")

type ChannelWithTimeout struct {
	ssh.Channel
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *ChannelWithTimeout) Read(data []byte) (int, error) {
	read, err := ioWithTimeout(c.Channel.Read, data, c.ReadTimeout)
	if err != nil && err == errTimeout {
		err = ErrTimeoutRead
	}
	return read, err
}

func (c *ChannelWithTimeout) Write(data []byte) (int, error) {
	written, err := ioWithTimeout(c.Channel.Write, data, c.WriteTimeout)
	if err != nil && err == errTimeout {
		err = ErrTimeoutWrite
	}
	return written, err
}

func ioWithTimeout(ioFunc func(data []byte) (int, error), data []byte, timeout time.Duration) (int, error) {
	timer := time.NewTimer(timeout)
	done := make(chan struct{})
	var (
		l   int
		err error
	)
	go func() {
		l, err = ioFunc(data)
		done <- struct{}{}
	}()
	select {
	case <-done:
	case <-timer.C:
		err = errTimeout
	}
	timer.Stop()
	return l, err
}
