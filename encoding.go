package gutil

import (
	"bytes"
)

//小端是以字节的角度来说的,低位字节在内存的低位就叫小端.
//但是在字节的内部还是从右往左一次升高的.
func EncodeULeb128(value uint64) []byte {
	//每次右移动7位高位补1,最后一次不补.
	//比如: 100110 0001110 1100101
	//变成: 0100110 10001110 11100101
	//变成: 0x26 0x8E 0xE5,对应的小端就是 0xE5 0x8E 0x26
	remaining := value >> 7
	var buf = new(bytes.Buffer)
	for remaining != 0 {
		buf.WriteByte(byte(value&0x7f | 0x80))
		value = remaining
		remaining >>= 7
	}
	buf.WriteByte(byte(value & 0x7f))
	return buf.Bytes()
}

//用read,因为是解码,要解很多次.
func DecodeULeb128(reader *bytes.Buffer) (uint64, int) {
	var (
		result uint64
		shift  uint64
		n      int
	)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			panic("could not parse LEB128 value")
		}
		//一次读byte,每次只拿7位,最高一位不要.
		//还要相应累加垫高7位
		result |= uint64((uint(b) & 0x7f) << shift)
		n++
		shift += 7
		//因为只有最后一个字节是0在最高位.
		//标示结束.
		if b&0x80 == 0 {
			break
		}
	}
	return result, n
}
