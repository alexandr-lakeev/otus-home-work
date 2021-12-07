package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *Telnet) Connect() error {
	var err error

	dialer := &net.Dialer{}

	t.ctx, t.cancel = context.WithTimeout(context.Background(), t.timeout)
	t.conn, err = dialer.DialContext(t.ctx, "tcp", t.address)
	if err != nil {
		return err
	}

	os.Stderr.Write([]byte(fmt.Sprintf("...Connected to %s\n", t.address)))

	return nil
}

func (t *Telnet) Close() error {
	t.cancel()
	return t.conn.Close()
}

func (t *Telnet) Send() error {
	scanner := bufio.NewScanner(t.in)
OUTER:
	for {
		select {
		case <-t.ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				os.Stderr.Write([]byte("...EOF\n"))
				break OUTER
			}
			_, err := t.conn.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				return err
			}
		}
	}
	t.Close()
	return nil
}

func (t *Telnet) Receive() error {
	scanner := bufio.NewScanner(t.conn)
OUTER:
	for {
		select {
		case <-t.ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				os.Stderr.Write([]byte("...Connection was closed by peer\n"))
				break OUTER
			}
			_, err := t.out.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				return err
			}
		}
	}
	t.Close()
	return nil
}
