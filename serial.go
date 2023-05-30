package modbus

import (
	"bytes"
	"fmt"
	"go.bug.st/serial"
	"sync"
	"time"
)

const (
	defaultSerialReadTimeOut = time.Second
	defaultBaudRate          = 9600
	// 开始3.5字符 地址8位 功能码8位 长度8位 CRC效验16位 结束3.5字符
	rtuMinByteLen int64 = 1 + 1 + 2
	// 串口传输一个字节位bit的长度 开始位固定为1 数据位8 奇偶校验位1 停止位1
	serialByteLen int64 = 1 + 8 + 1 + 1
)

type SerialPortTransporter struct {
	PortName string
	serial.Mode
	ReadTimeOut time.Duration
	port        serial.Port
	mu          sync.Mutex
}

func (t *SerialPortTransporter) Open() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.open()
}
func (t *SerialPortTransporter) open() error {
	port, err := serial.Open(t.PortName, &t.Mode)
	if err == nil {
		readTimeOut := defaultSerialReadTimeOut
		if t.ReadTimeOut > 0 {
			readTimeOut = t.ReadTimeOut
		}
		_ = port.SetReadTimeout(readTimeOut)
		t.port = port
	}
	return err
}
func (t *SerialPortTransporter) Connected() bool {
	return t.port != nil
}
func (t *SerialPortTransporter) Send(aduRequest ApplicationDataUnit) (aduResponse []byte, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.Connected() {
		err = t.open()
		if err != nil {
			return
		}
		if !t.Connected() {
			err = fmt.Errorf("serial could not open %s", t.PortName)
			return
		}
	}
	err = t.port.ResetInputBuffer()
	if err != nil {
		return
	}
	_, err = t.port.Write(aduRequest.GetData())
	if err != nil {
		return
	}
	sleep := t.calculateDelay(aduRequest)
	time.Sleep(sleep)
	var n int
	buf := bytes.NewBuffer([]byte{})
	for true {
		temp := make([]byte, rtuMaxSize*2)
		n, err = t.port.Read(temp)
		if err != nil {
			return
		}
		if n == 0 {
			break
		}
		buf.Write(temp[:n])
		if buf.Len() > rtuMinSize {
			break
		}
		time.Sleep(sleep)
	}
	aduResponse = buf.Bytes()
	return
}
func (t *SerialPortTransporter) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.close()
}

func (t *SerialPortTransporter) close() error {
	if !t.Connected() {
		return nil
	}
	err := t.port.Close()
	t.port = nil
	return err
}

func (t *SerialPortTransporter) calculateDelay(aduRequest ApplicationDataUnit) time.Duration {
	baudRate := t.BaudRate
	if baudRate == 0 {
		baudRate = defaultBaudRate
	}
	bitDelay := time.Second.Milliseconds() / int64(baudRate)
	length := 0
	switch aduRequest.GetFunctionCode() {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		count := aduRequest.Length()
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		count := aduRequest.Length()
		length += 1 + count*2
	case FuncCodeWriteSingleCoil,
		FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister,
		FuncCodeWriteMultipleRegisters:
		length += 4
	}
	byteDelay := bitDelay * serialByteLen
	return time.Duration(rtuMinByteLen*byteDelay + int64(length)*byteDelay)
}

func NewSerialTransporter(portName string) (t *SerialPortTransporter) {
	t = &SerialPortTransporter{
		PortName: portName,
	}
	return
}
