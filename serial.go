package modbus

import (
	"fmt"
	"go.bug.st/serial"
	"sync"
	"time"
)

const (
	defaultBaudRate = 9600
	// 开始3.5字符 地址8位 功能码8位 长度8位 CRC效验16位 结束3.5字符
	rtuMinByteLen int64 = 1 + 1 + 2
	// 串口传输一个字节位bit的长度 开始位固定为1 数据位8 奇偶校验位1 停止位1
	serialByteLen int64 = 1 + 8 + 1 + 1
)

type SerialPortTransporter struct {
	PortName string
	serial.Mode
	port serial.Port
	mu   sync.Mutex
}

func (t *SerialPortTransporter) Open() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.open()
}
func (t *SerialPortTransporter) open() error {
	port, err := serial.Open(t.PortName, &t.Mode)
	if err == nil {
		_ = port.SetReadTimeout(time.Second)
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
			return nil, err
		}
		if !t.Connected() {
			return nil, fmt.Errorf("serial could not open %s", t.PortName)
		}
	}
	_, err = t.port.Write(aduRequest.GetData())
	if err != nil {
		return nil, err
	}
	time.Sleep(t.calculateDelay(aduRequest))
	bl := len(aduRequest.GetCheckSumByte())
	max := rtuMaxSize
	// ascii 长度翻倍
	if bl == 1 {
		max = rtuMaxSize * 2
	}
	buf := make([]byte, max)
	n, err := t.port.Read(buf)
	if err != nil {
		return nil, err
	}
	data := buf[:n]
	return data, nil
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
