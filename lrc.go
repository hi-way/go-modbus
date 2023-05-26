package modbus

func LRC(bs []byte) uint8 {
	var sum uint8 = 0
	var b byte
	for _, b = range bs {
		sum += b
	}
	return uint8(-int8(sum))
}
