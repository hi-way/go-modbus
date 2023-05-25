package modbus

import (
	"encoding/hex"
)

const (
	//Bit 操作

	// FuncCodeReadDiscreteInputs 功能码:读离散量输入
	FuncCodeReadDiscreteInputs = 2
	// FuncCodeReadCoils 功能码:读线圈
	FuncCodeReadCoils = 1
	// FuncCodeWriteSingleCoil 功能码:写单个线圈
	FuncCodeWriteSingleCoil = 5
	// FuncCodeWriteMultipleCoils 功能码:写多个线圈
	FuncCodeWriteMultipleCoils = 15

	//16-bit 操作

	// FuncCodeReadInputRegisters 功能码:读输入寄存器
	FuncCodeReadInputRegisters = 4
	// FuncCodeReadHoldingRegisters 功能码:读保持寄存器
	FuncCodeReadHoldingRegisters = 3
	// FuncCodeWriteSingleRegister 功能码:写单个寄存器
	FuncCodeWriteSingleRegister = 6
	// FuncCodeWriteMultipleRegisters 功能码:写多个寄存器
	FuncCodeWriteMultipleRegisters = 16
	// FuncCodeReadWriteMultipleRegisters 功能码:读/写多个寄存器
	FuncCodeReadWriteMultipleRegisters = 23
)

// ProtocolDataUnit 协议数据单元
type ProtocolDataUnit interface {
	GetFunctionCode() (f byte)
	Length() (l int)
	GetData() (data []byte)
	ToHex() (h string)
}

// 协议数据单元
type protocolDataUnit struct {
	functionCode byte
	data         []byte
	length       int
}

func (pdu protocolDataUnit) GetFunctionCode() (f byte) {
	return pdu.functionCode
}

func (pdu protocolDataUnit) Length() (l int) {
	return pdu.length
}
func (pdu protocolDataUnit) GetData() (data []byte) {
	return pdu.data
}
func (pdu protocolDataUnit) ToHex() (h string) {
	return hex.EncodeToString(pdu.data)
}

// ApplicationDataUnit 应用数据单元
type ApplicationDataUnit interface {
	GetSlaveId() (slaveID byte)
	GetFunctionCode() (f byte)
	PDULength() (l int)
	GetPDU() (data []byte)
	ToHex() (h string)
	GetData() (data []byte)
	Length() (l int)
	GetCheckSum() (data uint16)
	GetCheckSumByte() (data []byte)
}

type applicationDataUnit struct {
	slaveID      byte
	pdu          ProtocolDataUnit
	data         []byte
	checkSum     uint16
	checkSumByte []byte
}

func (u applicationDataUnit) GetSlaveId() (slaveID byte) {
	return slaveID
}
func (u applicationDataUnit) GetFunctionCode() (f byte) {
	return u.pdu.GetFunctionCode()
}
func (u applicationDataUnit) GetPDU() (data []byte) {
	return u.pdu.GetData()
}
func (u applicationDataUnit) Length() (l int) {
	return len(u.data)
}

func (u applicationDataUnit) PDULength() (l int) {
	return u.pdu.Length()
}
func (u applicationDataUnit) ToHex() (h string) {
	return hex.EncodeToString(u.data)
}
func (u applicationDataUnit) GetData() (data []byte) {
	return u.data
}
func (u applicationDataUnit) GetCheckSum() (data uint16) {
	return u.checkSum
}
func (u applicationDataUnit) GetCheckSumByte() (data []byte) {
	return u.checkSumByte
}

// Packager 包解析器
type Packager interface {
	Encode(pdu protocolDataUnit) (adu ApplicationDataUnit, err error)
	Decode(results []byte) (adu ApplicationDataUnit, err error)
	Verify(aduRequest ApplicationDataUnit, aduResponse ApplicationDataUnit) (err error)
}

// Transporter 数据传输器
type Transporter interface {
	Open() error
	Connected() bool
	Send(aduRequest ApplicationDataUnit) (aduResponse []byte, err error)
	Close() error
}
