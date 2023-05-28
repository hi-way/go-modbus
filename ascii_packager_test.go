package modbus

import (
	"encoding/hex"
	"testing"
)

// 使用modbus slave 模拟串口3进行测试
func TestAsciiReadCoils(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewAsciiPackager(1)
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
func TestAsciiWriteSingleCoil(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewAsciiPackager(1)
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
func TestAsciiWriteMultipleCoils(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewAsciiPackager(1)
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
func TestAsciiReadDiscreteInputs(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewAsciiPackager(1)
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
func TestAsciiReadInputRegisters(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewAsciiPackager(1)
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
func TestAsciiReadHoldingRegisters(t *testing.T) {
	st := NewSerialTransporter("COM3")
	defer func() { _ = st.Close() }()
	pk := NewAsciiPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadHoldingRegisters(0, 10)
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
