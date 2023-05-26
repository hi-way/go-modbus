package modbus

import (
	"encoding/hex"
	"go.bug.st/serial"
	"testing"
)

// 使用modbus slave 模拟串口3进行测试
func TestRtuReadCoils(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadCoils(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestRtuWriteSingleCoil(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.WriteSingleCoil(0, true)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestRtuWriteMultipleCoils(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	coils := []bool{false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}
	request, results, err := c.WriteMultipleCoils(0, uint16(len(coils)), coils)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestRtuReadDiscreteInputs(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadDiscreteInputs(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestRtuReadInputRegisters(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewRtuPackager(1)
	c := NewClient(pk, ts)
	request, results, err := c.ReadInputRegisters(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestRtuReadHoldingRegisters(t *testing.T) {
	st := NewSerialTransporter("COM3")
	st.Mode = serial.Mode{BaudRate: defaultBaudRate}
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadHoldingRegisters(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
