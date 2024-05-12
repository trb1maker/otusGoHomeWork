package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
	_, err := io.Copy(c.con, c.in)
	return err
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.con)
	return err
}

func (c *client) Close() error {
	// Функции errors.Join нет в версии 1.19
	// return errors.Join(c.con.Close(), c.in.Close())
	return joinErrs(c.con.Close(), c.in.Close())
}

type errCollection []error

func (e errCollection) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}
	b := &strings.Builder{}
	for _, c := range e {
		b.WriteString(c.Error())
	}
	return b.String()
}

func (e errCollection) Unwrap() []error {
	return e
}

func joinErrs(errs ...error) error {
	err := make(errCollection, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			err = append(err, e)
		}
	}
	if len(err) != 0 {
		return err
	}
	return nil
}
