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
	return t.conn.Close()
}

func (t *Telnet) Send() error {
	buffer := make([]byte, 1024)
	for {
		n, err := t.in.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		n, err = t.conn.Write([]byte(buffer[:n]))
		if err != nil {
			return err
		}
	}
}

func (t *Telnet) Receive() error {
	buffer := make([]byte, 1024)
	for {
		n, err := t.conn.Read(buffer)
		if err != nil {
			t.in.Close()
			if err == io.EOF {
				os.Stderr.Write([]byte("...Connection was closed by peer\n"))
				return nil
			}
			return err
		}

		n, err = t.out.Write([]byte(string(buffer[:n])))
		if err != nil {
			return err
		}
	}
}
