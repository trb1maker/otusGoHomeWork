package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	con     io.ReadWriteCloser
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (c *client) Connect() error {
	var err error
	c.con, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("can't connect: %w", err)
	}
	return nil
}

func (c *client) Send() error {
	n, err := io.Copy(c.con, c.in)
	_ = n
	return err
}

func (c *client) Receive() error {
	n, err := io.Copy(c.out, c.con)
	_ = n
	return err
}

func (c *client) Close() error {
	return errors.Join(c.con.Close(), c.in.Close())
}
