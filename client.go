package modbus

import (
	"encoding/binary"
	"fmt"
)

// Client modbus客户端功能接口
type Client interface {
	//Bit 操作

	// ReadCoils 读线圈
	//地址:0~65535
	//范围:1~2000
	//功能码:1
	ReadCoils(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// ReadDiscreteInputs 读离散量输入
	//地址:0~65535
	//范围:1~2000
	//功能码:2
	ReadDiscreteInputs(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// WriteSingleCoil  写单个线圈
	//地址:0~65535
	//值:0或1
	//功能码:5
	WriteSingleCoil(address uint16, value bool) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// WriteMultipleCoils 写多个线圈
	//地址:0~65535
	//范围:1~1968
	//值:0或1
	//功能码:15
	WriteMultipleCoils(address, quantity uint16, value []bool) (request ApplicationDataUnit, results ApplicationDataUnit, err error)

	//16-bit 操作

	// ReadInputRegisters 读输入寄存器
	//地址:0~65535
	//范围:1~125
	//功能码:4
	ReadInputRegisters(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// ReadHoldingRegisters 读保持寄存器
	//地址:0~65535
	//范围:1~125
	//功能码:3
	ReadHoldingRegisters(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// WriteSingleRegister 写单个寄存器
	//地址:0~65535
	//值:0~65535
	//功能码:6
	WriteSingleRegister(address, value uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// WriteMultipleRegisters 写多个寄存器
	//地址:0~65535
	//范围:1~123
	//功能码:16
	WriteMultipleRegisters(address, quantity uint16, value []byte) (request ApplicationDataUnit, results ApplicationDataUnit, err error)
	// ReadWriteMultipleRegisters 读/写多个寄存器
	//地址:0~65535
	//读范围:1~125
	//写范围:1~121
	//功能码:23
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (request ApplicationDataUnit, results ApplicationDataUnit, err error)

	Send(request ApplicationDataUnit) (results ApplicationDataUnit, err error)
}

// 客户端 结构体
type client struct {
	packager    Packager
	transporter Transporter
}

func (c *client) ReadCoils(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v'", quantity, 1, 2000)
		return
	}
	data := dataBlock(address, quantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeReadCoils,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) ReadDiscreteInputs(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 2000)
		return
	}
	data := dataBlock(address, quantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeReadDiscreteInputs,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) WriteSingleCoil(address uint16, value bool) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if address > 65535 {
		err = fmt.Errorf("modbus: address '%v' must be between '%v' and '%v'", address, 0, 65535)
		return
	}
	// The requested ON/OFF state can only be 0xFF00 and 0x0000
	var coil uint16 = 0x0000
	if value {
		coil = 0xFF00
		return
	}
	data := dataBlock(address, coil)
	pdu := protocolDataUnit{
		functionCode: FuncCodeWriteSingleCoil,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) WriteMultipleCoils(address, quantity uint16, value []bool) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if quantity < 1 || quantity > 1968 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v'", quantity, 1, 1968)
		return
	}
	length := len(value)
	if int(quantity) != length {
		err = fmt.Errorf("modbus: quantity '%v' and value length '%v' are the same", quantity, length)
		return
	}
	coils := toBit(value)
	data := dataBlockSuffix(coils, address, quantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeWriteMultipleCoils,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) ReadInputRegisters(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}
	data := dataBlock(address, quantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeReadInputRegisters,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) ReadHoldingRegisters(address, quantity uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}
	data := dataBlock(address, quantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeReadHoldingRegisters,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) WriteSingleRegister(address, value uint16) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	data := dataBlock(address, value)
	pdu := protocolDataUnit{
		functionCode: FuncCodeWriteSingleRegister,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) WriteMultipleRegisters(address, quantity uint16, value []byte) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if quantity < 1 || quantity > 123 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 123)
		return
	}
	data := dataBlockSuffix(value, address, quantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeWriteMultipleRegisters,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}
func (c *client) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (request ApplicationDataUnit, results ApplicationDataUnit, err error) {
	if readQuantity < 1 || readQuantity > 125 {
		err = fmt.Errorf("modbus: quantity to read '%v' must be between '%v' and '%v',", readQuantity, 1, 125)
		return
	}
	if writeQuantity < 1 || writeQuantity > 121 {
		err = fmt.Errorf("modbus: quantity to write '%v' must be between '%v' and '%v',", writeQuantity, 1, 121)
		return
	}
	data := dataBlockSuffix(value, readAddress, readQuantity, writeAddress, writeQuantity)
	pdu := protocolDataUnit{
		functionCode: FuncCodeWriteMultipleRegisters,
		data:         data,
		length:       len(data),
	}
	request, err = c.packager.Encode(pdu)
	results, err = c.Send(request)
	return
}

var complementary = []byte{0x1, 0x2, 0x4, 0x8, 0x10, 0x20, 0x40, 0x80}

func toBit(value []bool) []byte {
	l := len(value)
	length := l / 8
	if l%8 > 0 {
		length++
	}
	bi := make([]byte, length)
	for i := 0; i < l; i++ {
		if value[i] {
			index := i / 8
			b := bi[index]
			s := i % 8
			c := complementary[s]
			bi[index] = b | c
		}
	}
	return bi
}
func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return data
}
func dataBlockSuffix(suffix []byte, value ...uint16) []byte {
	length := 2 * len(value)
	data := make([]byte, length+1+len(suffix))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	data[length] = uint8(len(suffix))
	copy(data[length+1:], suffix)
	return data
}
func (c *client) Send(request ApplicationDataUnit) (results ApplicationDataUnit, err error) {
	bys, err := c.transporter.Send(request)
	if err != nil {
		return
	}
	results, err = c.packager.Decode(bys)
	if err != nil {
		return
	}
	err = c.packager.Verify(request, results)
	return
}
func NewClient(packager Packager, transporter Transporter) (c Client) {
	c = &client{
		packager:    packager,
		transporter: transporter,
	}
	return
}
