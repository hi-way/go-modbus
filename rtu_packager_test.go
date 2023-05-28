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
		t.Log("request", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("request", hex.EncodeToString(request.GetData()))
	t.Log("results", hex.EncodeToString(results.GetData()))
	t.Logf("results length %d ", results.GetPDU().Length())
	t.Logf("results data %s", hex.EncodeToString(results.GetPDU().GetData()))
}
func TestRtuWriteSingleCoil(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.WriteSingleCoil(0, true)
	if err != nil {
		t.Log("request", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("request", hex.EncodeToString(request.GetData()))
	t.Log("results", hex.EncodeToString(results.GetData()))
	t.Logf("results length %d ", results.GetPDU().Length())
	t.Logf("results data %s", hex.EncodeToString(results.GetPDU().GetData()))
}
func TestRtuWriteMultipleCoils(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	coils := []bool{false, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}
	request, results, err := c.WriteMultipleCoils(0, uint16(len(coils)), coils)
	if err != nil {
		t.Log("request", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("request", hex.EncodeToString(request.GetData()))
	t.Log("results", hex.EncodeToString(results.GetData()))
	t.Logf("results length %d ", results.GetPDU().Length())
	t.Logf("results data %s", hex.EncodeToString(results.GetPDU().GetData()))
}
func TestRtuReadDiscreteInputs(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadDiscreteInputs(0, 10)
	if err != nil {
		t.Log("request", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("request", hex.EncodeToString(request.GetData()))
	t.Log("results", hex.EncodeToString(results.GetData()))
	t.Logf("results length %d ", results.GetPDU().Length())
	t.Logf("results data %s", hex.EncodeToString(results.GetPDU().GetData()))
}
func TestRtuReadInputRegisters(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewRtuPackager(1)
	c := NewClient(pk, ts)
	request, results, err := c.ReadInputRegisters(0, 10)
	if err != nil {
		t.Log("request", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("request", hex.EncodeToString(request.GetData()))
	t.Log("results", hex.EncodeToString(results.GetData()))
	t.Logf("results length %d ", results.GetPDU().Length())
	t.Logf("results data %s", hex.EncodeToString(results.GetPDU().GetData()))
}
func TestRtuReadHoldingRegisters(t *testing.T) {
	st := NewSerialTransporter("COM3")
	st.Mode = serial.Mode{BaudRate: defaultBaudRate}
	defer func() { _ = st.Close() }()
	pk := NewRtuPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadHoldingRegisters(0, 1)
	if err != nil {
		t.Log("request", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("request", hex.EncodeToString(request.GetData()))
	t.Log("results", hex.EncodeToString(results.GetData()))
	t.Logf("results length %d ", results.GetPDU().Length())
	t.Logf("results data %s", hex.EncodeToString(results.GetPDU().GetData()))
}
