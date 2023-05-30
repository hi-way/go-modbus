package modbus

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	defaultTcpConnectTimeout = 1 * time.Second
	defaultTcpReadTimeout    = 1 * time.Second
	defaultTcpWriteTimeout   = 1 * time.Second
	defaultTcpKeepAlive      = 30 * time.Second
)

type TcpTransporter struct {
	Address        string
	ConnectTimeout time.Duration
	KeepAlive      time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	mu             sync.Mutex
	conn           net.Conn
}

func (mb *TcpTransporter) Send(aduRequest ApplicationDataUnit) (aduResponse []byte, err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if !mb.Connected() {
		err = mb.connect()
		if err != nil {
			return
		}
	}
	tcpWriteTimeout := defaultTcpWriteTimeout
	if mb.WriteTimeout > 0 {
		tcpWriteTimeout = mb.WriteTimeout
	}
	err = mb.conn.SetWriteDeadline(time.Now().Add(tcpWriteTimeout))
	if err != nil {
		_ = mb.close()
		return
	}
	_, err = mb.conn.Write(aduRequest.GetData())
	if err != nil {
		_ = mb.close()
		return
	}
	tcpReadTimeout := defaultTcpReadTimeout
	if mb.ReadTimeout > 0 {
		tcpReadTimeout = mb.ReadTimeout
	}
	err = mb.conn.SetReadDeadline(time.Now().Add(tcpReadTimeout))
	if err != nil {
		_ = mb.close()
		return
	}
	temp := make([]byte, tcpMaxSize*2)
	rl, err := mb.conn.Read(temp)
	if err != nil && rl == 0 {
		_ = mb.close()
		return
	}
	if rl <= 0 {
		err = fmt.Errorf("modbus: Read  data is  empty")
		return
	}
	aduResponse = make([]byte, rl)
	copy(aduResponse, temp[:rl])
	return
}
func (mb *TcpTransporter) Connected() bool {
	return mb.conn != nil
}
func (mb *TcpTransporter) Open() error {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return mb.connect()
}

func (mb *TcpTransporter) connect() error {
	if mb.conn == nil {
		tcpConnectTimeout := defaultTcpConnectTimeout
		tcpKeepAlive := defaultTcpKeepAlive
		if mb.ConnectTimeout > 0 {
			tcpConnectTimeout = mb.ConnectTimeout
		}
		if mb.KeepAlive > 0 {
			tcpKeepAlive = mb.KeepAlive
		}
		dialer := net.Dialer{Timeout: tcpConnectTimeout, KeepAlive: tcpKeepAlive}
		conn, err := dialer.Dial("tcp", mb.Address)
		if err != nil {
			return err
		}
		mb.conn = conn
	}
	return nil
}

func (mb *TcpTransporter) Close() error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	return mb.close()
}

func (mb *TcpTransporter) close() (err error) {
	if mb.conn != nil {
		err = mb.conn.Close()
		mb.conn = nil
	}
	return
}

func NewTcpTransporter(address string) (t *TcpTransporter) {
	t = &TcpTransporter{
		Address: address,
	}
	return
}
