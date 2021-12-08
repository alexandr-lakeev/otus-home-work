package main

import (
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

	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	os.Stderr.Write([]byte(fmt.Sprintf("...Connected to %s\n", t.address)))
	return nil
}

func (t *Telnet) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *Telnet) Send() error {
	if _, err := io.Copy(t.conn, t.in); err != nil {
		os.Stderr.Write([]byte(err.Error()))

		return err
	}

	os.Stderr.Write([]byte("...EOF"))
	return nil
}

func (t *Telnet) Receive() error {
	if _, err := io.Copy(t.out, t.conn); err != nil {
		return err
	}

	os.Stderr.Write([]byte("...Connection was closed by peer"))
	return nil
}
