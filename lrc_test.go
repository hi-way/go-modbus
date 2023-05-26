package modbus

import (
	"encoding/hex"
	"testing"
)

func TestLRC(t *testing.T) {
	s := "3031303330303030303030314642"
	bs, _ := hex.DecodeString(s)
	lrc := LRC(bs)
	t.Log(lrc)
}
