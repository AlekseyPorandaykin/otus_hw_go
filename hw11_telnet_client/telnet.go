package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

var (
	ErrEOF               = errors.New("end read file")
	ErrReadConnectClosed = errors.New("connect for read closed")
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type TCPTelnetClient struct {
	address string
	conn    net.Conn
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *TCPTelnetClient) Connect() error {
	conn, errC := net.DialTimeout("tcp", t.address, t.timeout)
	if errC != nil {
		return errC
	}
	log.Printf("...Connected %s \n", t.address)
	t.conn = conn

	return nil
}

func (t *TCPTelnetClient) Close() error {
	if t.conn != nil {
		if errConn := t.conn.Close(); errConn != nil {
			return errConn
		}
	}
	if errI := t.in.Close(); errI != nil {
		return errI
	}
	log.Println("...Connection was closed by peer")
	return nil
}

func (t *TCPTelnetClient) Send() error {
	return resendMessages(t.in, t.conn)
}

func (t *TCPTelnetClient) Receive() error {
	return resendMessages(t.conn, t.out)
}

func resendMessages(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		errR := s.Err()
		message := s.Text()
		if errors.Is(errR, io.EOF) {
			log.Println("...EOF")
			return ErrEOF
		}
		if errR != nil {
			if errors.Is(errR, net.ErrClosed) {
				return ErrReadConnectClosed
			}
			return errR
		}
		_, errW := w.Write([]byte(message + "\n"))
		if errW != nil {
			return errW
		}
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TCPTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
