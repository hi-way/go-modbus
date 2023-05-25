package modbus

import (
	"encoding/hex"
	"testing"
)

// 使用modbus slave 模拟串口3进行测试
func TestReadCoils(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewRtuPackager(1)
	c := NewClient(pk, ts)
	request, results, err := c.ReadCoils(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestWriteSingleCoil(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewRtuPackager(1)
	c := NewClient(pk, ts)
	request, results, err := c.WriteSingleCoil(0, true)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
func TestWriteMultipleCoils(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewRtuPackager(1)
	c := NewClient(pk, ts)
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
func TestReadDiscreteInputs(t *testing.T) {
	ts := NewSerialTransporter("COM3")
	pk := NewRtuPackager(1)
	c := NewClient(pk, ts)
	request, results, err := c.ReadDiscreteInputs(0, 10)
	if err != nil {
		t.Log("Write", hex.EncodeToString(request.GetData()))
		t.Error(err)
		t.FailNow()
	}
	t.Log("Write", hex.EncodeToString(request.GetData()))
	t.Log(hex.EncodeToString(results.GetData()))
}
