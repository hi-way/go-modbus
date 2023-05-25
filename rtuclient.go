package modbus

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

const (
	rtuMinSize = 4
	rtuMaxSize = 256
)

// rtuPackager  rtu包解析器
type rtuPackager struct {
	slaveID byte
}

func (r *rtuPackager) Encode(pdu protocolDataUnit) (adu ApplicationDataUnit, err error) {
	length := len(pdu.GetData()) + 4
	if length > rtuMaxSize {
		err = fmt.Errorf("modbus: length of data '%v' must not be bigger than '%v'", length, rtuMaxSize)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(r.slaveID)
	buf.WriteByte(pdu.GetFunctionCode())
	buf.Write(pdu.GetData())
	checksum := CRC16(buf.Bytes())
	checksumByte := CRC16ToBytes(checksum)
	buf.Write(checksumByte)
	adu = applicationDataUnit{
		slaveID:      r.slaveID,
		pdu:          pdu,
		checkSum:     checksum,
		checkSumByte: checksumByte,
		data:         buf.Bytes(),
	}
	return
}
func (r *rtuPackager) Decode(results []byte) (adu ApplicationDataUnit, err error) {
	length := len(results)
	if length > rtuMaxSize {
		err = fmt.Errorf("modbus: response data size '%v' exceeds the maximum limit of '%v'", length, rtuMaxSize)
		return
	}
	if length < rtuMinSize {
		err = fmt.Errorf("modbus: response data size '%v' less than maximum limit of '%v'", length, rtuMinSize)
		return
	}
	slaveID := results[0]
	functionCode := results[1]
	pduData := results[1 : length-2]
	pduLength := results[2]
	pdu := protocolDataUnit{
		functionCode: functionCode,
		data:         pduData,
		length:       int(pduLength),
	}
	checkSumByte := results[length-2:]
	adu = applicationDataUnit{
		slaveID:      slaveID,
		pdu:          pdu,
		checkSumByte: checkSumByte,
		checkSum:     CRC16ToUint(checkSumByte),
		data:         results,
	}
	return
}
func (r *rtuPackager) Verify(aduRequest ApplicationDataUnit, aduResponse ApplicationDataUnit) (err error) {
	if aduRequest.GetSlaveId() != aduResponse.GetSlaveId() {
		err = fmt.Errorf("modbus: aduRequest  slaveId '%v' and aduResponse slaveId '%v' are inconsistent", aduRequest.GetSlaveId(), aduResponse.GetSlaveId())
		return
	}
	if aduRequest.GetFunctionCode() != aduResponse.GetFunctionCode() {
		err = fmt.Errorf("modbus: aduRequest  functionCode '%v' and aduResponse functionCode '%v' are inconsistent", aduRequest.GetFunctionCode(), aduResponse.GetFunctionCode())
		return
	}
	data := aduResponse.GetData()
	length := len(data)
	checksum := CRC16(aduResponse.GetData()[:length-2])
	if checksum != aduResponse.GetCheckSum() {
		sumByte := make([]byte, 2)
		sumByte[0] = byte(checksum >> 8)
		sumByte[1] = byte(checksum)
		err = fmt.Errorf("modbus: crc validation failed source:'%v' reality:'%v' ", hex.EncodeToString(aduResponse.GetCheckSumByte()), hex.EncodeToString(sumByte))
		return
	}
	return
}

func NewRtuPackager(slaveID byte) (p Packager) {
	p = &rtuPackager{slaveID: slaveID}
	return
}

func NewSerialTransporter(portName string) (t Transporter) {
	t = &serialPortTransporter{
		portName: portName,
	}
	return
}
