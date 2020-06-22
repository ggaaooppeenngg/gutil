package util

import (
	"bytes"
	"testing"
)

func TestLEB128(t *testing.T) {
	bs := EncodeULeb128(624485)
	num, _ := DecodeULeb128(bytes.NewBuffer(bs))
	if !bytes.Equal(bs, []byte{0xE5, 0x8E, 0x26}) {
		t.Fatal("encode error")
	}
	if num != 624485 {
		t.Fatal("decode error")
	}
}
