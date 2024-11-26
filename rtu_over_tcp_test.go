package modbus

import (
	"encoding/hex"
	"testing"
)

// 使用modbus slave 模拟进行测试
func TestRtuTcpReadCoils(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
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
func TestRtuTcpWriteSingleCoil(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
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
func TestRtuTcpWriteMultipleCoils(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
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
func TestRtuTcpReadDiscreteInputs(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
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
func TestRtuTcpReadInputRegisters(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadInputRegisters(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestRtuTcpReadHoldingRegisters(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
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
