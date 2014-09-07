//
// Copyright 2014 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"math"
)

func GenerateRandomBytes(size uint) []byte {
	size64 := size / uint(8)
	buf := new(bytes.Buffer)
	var r64 int64
	var i uint
	for i = uint(0); i < size64; i++ {
		r64 = rand.Int63()
		binary.Write(buf, binary.BigEndian, &r64)
	}
	for i = i * uint(8); i < size; i++ {
		buf.WriteByte(byte(rand.Int()))
	}
	return buf.Bytes()
}

func ReadUint8(b []byte) (v uint8) {
	v = uint8(b[0])
	return v
}

func ReadUint16(b []byte) (v uint16) {
	v = uint16(b[1]) | uint16(b[0])<<8
	return v
}

func ReadUint24(b []byte) (v uint32) {
	v = uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
	return v
}

func ReadUint32(b []byte) (v uint32) {
	v = uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	return v
}


func ReadFloat64(b []byte) (v float64) {
	v64 := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	v = math.Float64frombits(v64)

	return v
}