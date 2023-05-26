package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

const (
	tcpProtocolIdentifier uint16 = 0x0000
	// Modbus Application Protocol
	tcpHeaderSize = 7
	tcpMaxSize    = 256
)

// tcpPackager  tcp包解析器
type tcpPackager struct {
	transactionID uint16
	slaveID       byte
	mux           sync.Mutex
}

func (p *tcpPackager) transaction() uint16 {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.transactionID++
	return p.transactionID
}

func (p *tcpPackager) Encode(pdu protocolDataUnit) (adu ApplicationDataUnit, err error) {
	aduLength := tcpHeaderSize + 1 + len(pdu.GetData())
	if aduLength > tcpMaxSize {
		err = fmt.Errorf("modbus: length of data '%v' must not be bigger than '%v'", aduLength, tcpMaxSize)
		return
	}
	bs := bytes.NewBuffer([]byte{})
	// Transaction identifier
	_ = binary.Write(bs, binary.BigEndian, p.transaction())
	// Protocol identifier
	_ = binary.Write(bs, binary.BigEndian, tcpProtocolIdentifier)
	// Length = sizeof(SlaveID) + sizeof(FunctionCode) + Data
	length := uint16(1 + 1 + len(pdu.GetData()))
	_ = binary.Write(bs, binary.BigEndian, length)
	// slaveID
	bs.WriteByte(p.slaveID)
	// PDU
	bs.WriteByte(pdu.GetFunctionCode())
	bs.Write(pdu.GetData())
	adu = applicationDataUnit{
		slaveID:      p.slaveID,
		pdu:          pdu,
		checkSum:     0,
		checkSumByte: []byte{0},
		data:         bs.Bytes(),
	}
	return
}
func (p *tcpPackager) Decode(results []byte) (adu ApplicationDataUnit, err error) {
	allLength := len(results)
	if allLength > tcpMaxSize {
		err = fmt.Errorf("modbus: response data size '%v' exceeds the maximum limit of '%v'", allLength, tcpMaxSize)
		return
	}
	if allLength < tcpHeaderSize {
		err = fmt.Errorf("modbus: response data size '%v' less than maximum limit of '%v'", allLength, tcpHeaderSize)
		return
	}
	// Read length value in the header
	length := binary.BigEndian.Uint16(results[4:])
	pduLength := allLength - tcpHeaderSize + 1
	if pduLength <= 0 || pduLength != int(length) {
		err = fmt.Errorf("modbus: length in response '%v' does not match pdu data length '%v'", length-1, pduLength)
		return
	}
	functionCode := results[tcpHeaderSize]
	pduData := results[tcpHeaderSize+1:]
	slaveID := results[tcpHeaderSize-1]
	pdu := protocolDataUnit{
		functionCode: functionCode,
		data:         pduData,
		length:       pduLength,
	}
	adu = applicationDataUnit{
		slaveID:      slaveID,
		pdu:          pdu,
		checkSumByte: []byte{0},
		checkSum:     0,
		data:         results,
	}
	return
}
func (p *tcpPackager) Verify(aduRequest ApplicationDataUnit, aduResponse ApplicationDataUnit) (err error) {
	if aduRequest.GetSlaveId() != aduResponse.GetSlaveId() {
		err = fmt.Errorf("modbus: aduRequest  slaveId '%v' and aduResponse slaveId '%v' are inconsistent", aduRequest.GetSlaveId(), aduResponse.GetSlaveId())
		return
	}
	if aduRequest.GetFunctionCode() != aduResponse.GetFunctionCode() {
		if aduResponse.GetFunctionCode() == aduRequest.GetFunctionCode()+0x80 {
			err = fmt.Errorf("modbus: error   errorCode '%v'", aduResponse.GetFunctionCode())
			return
		}
		err = fmt.Errorf("modbus: aduRequest  functionCode '%v' and aduResponse functionCode '%v' are inconsistent", aduRequest.GetFunctionCode(), aduResponse.GetFunctionCode())
		return
	}
	requestData := aduRequest.GetData()
	responseData := aduResponse.GetData()
	requestTransaction := binary.BigEndian.Uint16(requestData)
	responseTransaction := binary.BigEndian.Uint16(responseData)
	if requestTransaction != responseTransaction {
		err = fmt.Errorf("modbus: response transaction id '%v' does not match request '%v'", responseTransaction, requestTransaction)
		return
	}
	if requestData[6] != responseData[6] {
		err = fmt.Errorf("modbus: response unit id '%v' does not match request '%v'", responseData[6], requestData[6])
		return
	}
	return
}

func NewTcpPackager(slaveID byte) (p Packager) {
	p = &tcpPackager{slaveID: slaveID}
	return
}
