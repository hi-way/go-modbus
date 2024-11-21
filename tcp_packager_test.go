package modbus

import (
	"encoding/hex"
	"sync"
	"testing"
)

// 使用modbus slave 模拟进行测试
func TestTcpReadCoils(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewTcpPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadCoils(0, 2000)
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
func TestTcpWriteSingleCoil(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewTcpPackager(1)
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
func TestTcpWriteMultipleCoils(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewTcpPackager(1)
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
func TestTcpReadDiscreteInputs(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewTcpPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadDiscreteInputs(0, 2000)
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
func TestTcpReadInputRegisters(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewTcpPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadInputRegisters(0, 125)
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
func TestTcpReadHoldingRegisters(t *testing.T) {
	st := NewTcpTransporter("127.0.0.1:502")
	defer func() { _ = st.Close() }()
	pk := NewTcpPackager(1)
	c := NewClient(pk, st)
	request, results, err := c.ReadHoldingRegisters(0, 125)
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

func TestTransactionID(t *testing.T) {
	pk := tcpPackager{
		slaveID: 1,
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			t.Logf("tcpPackager TransactionID %d", pk.transaction())
			wg.Done()
		}()
	}

	wg.Wait()
}
