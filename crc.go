package modbus

import (
	"bytes"
	"encoding/binary"
	"sync"
)

// Cyclical Redundancy Checking.
var (
	once     sync.Once
	crtTable []uint16
)

func CRC16(bs []byte) uint16 {
	once.Do(initCrcTable)
	val := uint16(0xFFFF)
	for _, v := range bs {
		val = (val >> 8) ^ crtTable[(val^uint16(v))&0x00FF]
	}
	return val
}
func CRC16ToBytes(crc uint16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.LittleEndian, crc)
	return bytesBuffer.Bytes()
}

func CRC16ToUint(bs []byte) uint16 {
	return binary.LittleEndian.Uint16(bs)
}

// initCrcTable 初始化表.
func initCrcTable() {
	crcPoly16 := uint16(0xa001)
	crtTable = make([]uint16, 256)
	for i := uint16(0); i < 256; i++ {
		crc := uint16(0)
		b := i
		for j := uint16(0); j < 8; j++ {
			if ((crc ^ b) & 0x0001) > 0 {
				crc = (crc >> 1) ^ crcPoly16
			} else {
				crc >>= 1
			}
			b >>= 1
		}
		crtTable[i] = crc
	}
}
