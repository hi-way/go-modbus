package modbus

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	defaultTcpConnectTimeOut = 1 * time.Second
	defaultTcpReadTimeOut    = 1 * time.Second
	defaultTcpWriteTimeOut   = 1 * time.Second
	defaultTcpKeepAlive      = 30 * time.Second
)

type TcpTransporter struct {
	Address        string
	ConnectTimeOut time.Duration
	KeepAlive      time.Duration
	ReadTimeOut    time.Duration
	WriteTimeOut   time.Duration
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
	tcpWriteTimeOut := defaultTcpWriteTimeOut
	if mb.WriteTimeOut > 0 {
		tcpWriteTimeOut = mb.WriteTimeOut
	}
	err = mb.conn.SetWriteDeadline(time.Now().Add(tcpWriteTimeOut))
	if err != nil {
		_ = mb.close()
		return
	}
	_, err = mb.conn.Write(aduRequest.GetData())
	if err != nil {
		_ = mb.close()
		return
	}
	tcpReadTimeOut := defaultTcpReadTimeOut
	if mb.ReadTimeOut > 0 {
		tcpReadTimeOut = mb.ReadTimeOut
	}
	err = mb.conn.SetReadDeadline(time.Now().Add(tcpReadTimeOut))
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
		tcpConnectTimeOut := defaultTcpConnectTimeOut
		tcpKeepAlive := defaultTcpKeepAlive
		if mb.ConnectTimeOut > 0 {
			tcpConnectTimeOut = mb.ConnectTimeOut
		}
		if mb.KeepAlive > 0 {
			tcpKeepAlive = mb.KeepAlive
		}
		dialer := net.Dialer{Timeout: tcpConnectTimeOut, KeepAlive: tcpKeepAlive}
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
