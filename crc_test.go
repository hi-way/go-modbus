package modbus

import (
	"encoding/hex"
	"testing"
)

func TestCRC(t *testing.T) {
	crcr := "bc0d"
	str := "01010000000A"
	bs, _ := hex.DecodeString(str)
	crc := CRC16(bs)
	crcs := hex.EncodeToString(CRC16ToBytes(crc))
	t.Logf("'%s' crc '%d' hex '%s'", str, crc, crcs)
	if crcr != crcs {
		t.FailNow()
	}
}

func TestBitShift(t *testing.T) {
	value := []bool{true, true, true, false, false, false, true, true, true, false}
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
	t.Logf("bit %s", hex.EncodeToString(bi))
}
