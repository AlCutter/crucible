// crucible
// One-Time-Programmable (OTP) fusing tool
//
// Copyright (c) F-Secure Corporation
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package util

import (
	"fmt"
	"math/big"
	"math/bits"
)

// Pad4 pads a byte array to ensure that it always represents one or more
// 4-byte register.
func Pad4(val []byte) (res []byte) {
	numRegisters := 1 + len(val)/4

	// normalize
	if len(val)%4 == 0 {
		numRegisters -= 1
	}

	pad := numRegisters*4 - len(val)

	for i := 0; i < pad; i++ {
		val = append(val, 0x00)
	}

	return val
}

// PadBigInt pads a big.Int value to account for the fact that big.Bytes()
// returns the absolute value, therefore leading 0x00 bytes are not returned
// and 0x00 values are empty.
func PadBigInt(val *big.Int, size int) (res []byte) {
	numBytes := 1 + size/8

	// normalize
	if size%8 == 0 {
		numBytes -= 1
	}

	pad := numBytes - len(val.Bytes())
	res = val.Bytes()

	for i := 0; i < pad; i++ {
		res = append([]byte{0x00}, res...)
	}

	return
}

// SwitchEndianness reverses a byte array to switch between big <> little
// endianess.
func SwitchEndianness(val []byte) []byte {
	for i := len(val)/2 - 1; i >= 0; i-- {
		rev := len(val) - 1 - i
		val[i], val[rev] = val[rev], val[i]
	}

	return val
}

// ConvertReadValue converts a little-endian byte array by shifting it
// according to its register offset and size and converting it to big-endian.
func ConvertReadValue(off int, size int, val []byte) (res []byte) {
	// little-endian > big-endian
	res = SwitchEndianness(val)

	v := new(big.Int)
	v.SetBytes(res)
	v.Rsh(v, uint(off))

	// get only the bits that we care about
	mask := big.NewInt((1 << size) - 1)
	v.And(v, mask)

	res = PadBigInt(v, size)

	return
}

// ConvertWriteValue converts a big-endian byte array by shifting it according
// to its register offset and size and converting it to little-endian.
func ConvertWriteValue(off int, size int, val []byte) (res []byte, err error) {
	bitLen := bits.Len(uint(val[0])) + (len(val)-1)*8

	if bitLen > size {
		err = fmt.Errorf("value bit size %d exceeds %d", bitLen, size)
		return
	}

	v := new(big.Int)
	v.SetBytes(val)
	v.Lsh(v, uint(off))

	res = PadBigInt(v, size)
	// big-endian > little-endian
	res = SwitchEndianness(res)

	return
}
