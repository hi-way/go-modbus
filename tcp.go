package modbus

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	defaultTimeOut   = 1 * time.Second
	defaultKeepAlive = 30 * time.Second
)

type TcpTransporter struct {
	Address   string
	TimeOut   time.Duration
	KeepAlive time.Duration
	mu        sync.Mutex
	conn      net.Conn
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
	err = mb.conn.SetWriteDeadline(time.Now().Add(time.Second))
	if err != nil {
		_ = mb.close()
		return
	}
	_, err = mb.conn.Write(aduRequest.GetData())
	if err != nil {
		_ = mb.close()
		return
	}
	timeOut := defaultTimeOut
	if mb.TimeOut > 0 {
		timeOut = mb.TimeOut
	}
	err = mb.conn.SetReadDeadline(time.Now().Add(timeOut))
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
		timeOut := defaultTimeOut
		keepAlive := defaultKeepAlive
		if mb.TimeOut > 0 {
			timeOut = mb.TimeOut
		}
		if mb.KeepAlive > 0 {
			keepAlive = mb.KeepAlive
		}
		dialer := net.Dialer{Timeout: timeOut, KeepAlive: keepAlive}
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
