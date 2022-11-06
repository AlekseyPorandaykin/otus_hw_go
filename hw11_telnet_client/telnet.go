package main

import (
	"bufio"
	"context"
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
	io.Closer
	Send() error
	Receive() error
}

type TcpTelnetClient struct {
	address string
	conn    net.Conn
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *TcpTelnetClient) Connect() error {
	conn, errC := net.DialTimeout("tcp", t.address, t.timeout)
	if errC != nil {
		return errC
	}
	log.Printf("...Connected %s \n", t.address)
	t.conn = conn

	return nil
}

func (t *TcpTelnetClient) Close() error {
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

func (t *TcpTelnetClient) Send() error {
	return resendMessages(t.in, t.conn)
}

func (t *TcpTelnetClient) Receive() error {
	return resendMessages(t.conn, t.out)
}

func ReceiveFromServer(ctx context.Context, client TelnetClient) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Receive(); err != nil {
					log.Println(err)
					cancel()
					return
				}
			}
		}

	}()
}

func SendToServer(ctx context.Context, client TelnetClient) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Send(); err != nil {
					log.Println(err)
					cancel()
					return
				}
			}
		}
	}()
}

func resendMessages(r io.Reader, w io.Writer) error {
	message, errR := bufio.NewReader(r).ReadString('\n')
	if errR == io.EOF {
		log.Println("...EOF")
		return ErrEOF
	}
	if errR != nil {
		if errors.Is(errR, net.ErrClosed) {
			return ErrReadConnectClosed
		}
		return errR
	}
	_, errW := w.Write([]byte(message))
	if errW != nil {
		return errW
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TcpTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
