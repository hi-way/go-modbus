package modbus

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	asciiMaxSize = rtuMaxSize * 2
	asciiMinSize = rtuMinSize * 2
	asciiStart   = ":"
	asciiEnd     = "\r\n"
)

type asciiPackager struct {
	slaveID byte
}

func (p *asciiPackager) Encode(pdu protocolDataUnit) (adu ApplicationDataUnit, err error) {
	data := make([]byte, 3+len(pdu.GetData()))
	data[0] = p.slaveID
	data[1] = pdu.GetFunctionCode()
	copy(data[2:], pdu.GetData())
	checkSum := LRC(data)
	checkSumByte := []byte{checkSum}
	buf := bytes.NewBuffer([]byte{})
	//start
	_, _ = buf.WriteString(asciiStart)
	_, _ = buf.WriteString(strings.ToUpper(hex.EncodeToString(data)))
	_, _ = buf.WriteString(strings.ToUpper(hex.EncodeToString(checkSumByte)))
	_, _ = buf.WriteString(asciiEnd)
	adu = applicationDataUnit{
		slaveID:      p.slaveID,
		pdu:          pdu,
		checkSum:     uint16(checkSum),
		checkSumByte: checkSumByte,
		data:         buf.Bytes(),
	}
	return
}
func (p *asciiPackager) Decode(results []byte) (adu ApplicationDataUnit, err error) {
	length := len(results)
	if length > asciiMaxSize {
		err = fmt.Errorf("modbus: response data size '%v' exceeds the maximum limit of '%v'", length, asciiMaxSize)
		return
	}
	if length < asciiMinSize {
		err = fmt.Errorf("modbus: response data size '%v' less than maximum limit of '%v'", length, asciiMinSize)
		return
	}
	slaveIDBs, err := hex.DecodeString(string(results[1:3]))
	if err != nil {
		return
	}
	slaveID := slaveIDBs[0]
	functionCodeBs, err := hex.DecodeString(string(results[3:5]))
	if err != nil {
		return
	}
	functionCode := functionCodeBs[0]
	pduLengthBs, err := hex.DecodeString(string(results[5:7]))
	if err != nil {
		return
	}
	pduLength := pduLengthBs[0]
	pduData := results[7 : length-4]
	pduData, err = hex.DecodeString(string(pduData))
	if err != nil {
		return
	}
	pdu := protocolDataUnit{
		functionCode: functionCode,
		data:         pduData,
		length:       int(pduLength),
	}
	checkSum := results[length-4 : length-2]
	checkSum, err = hex.DecodeString(string(checkSum))
	if err != nil {
		return
	}
	adu = applicationDataUnit{
		slaveID:      slaveID,
		pdu:          pdu,
		checkSumByte: checkSum,
		checkSum:     uint16(checkSum[0]),
		data:         results,
	}
	return
}
func (p *asciiPackager) Verify(aduRequest ApplicationDataUnit, aduResponse ApplicationDataUnit) (err error) {
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
	data := aduResponse.GetData()
	data, err = hex.DecodeString(string(data[1 : len(data)-4]))
	if err != nil {
		return
	}
	checksum := LRC(data)
	if uint16(checksum) != aduResponse.GetCheckSum() {
		err = fmt.Errorf("modbus: lrc validation failed source:'%v' reality:'%v' ", hex.EncodeToString(aduResponse.GetCheckSumByte()), hex.EncodeToString([]byte{checksum}))
		return
	}
	return
}

func NewAsciiPackager(slaveID byte) (p Packager) {
	p = &asciiPackager{slaveID: slaveID}
	return
}
