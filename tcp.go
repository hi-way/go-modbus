package modbus

import (
	"fmt"
	"net"
	"sync"
	"time"
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
	err = mb.conn.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		_ = mb.close()
		return
	}
	temp := make([]byte, tcpMaxSize)
	rl, err := mb.conn.Read(temp)
	if err != nil {
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
		dialer := net.Dialer{Timeout: mb.TimeOut, KeepAlive: mb.KeepAlive}
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
