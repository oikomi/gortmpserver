//
// Copyright 2014-2099 Hong Miao. All Rights Reserved.
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

import(
	"math/rand"
	"bytes"
	"encoding/binary"
	//"github.com/golang/glog"
)

const (
	UBigEndian     = 0
	ULittleEndian  = 1
)

func Byte42Uint32(data []byte, endian int) uint32 {
	var i uint32
	
	if UBigEndian == endian {
		i = uint32(uint32(data[3]) + uint32(data[2])<<8 + uint32(data[1])<<16 + uint32(data[0])<<24)
	}
	
	if ULittleEndian == endian {
		i = uint32(uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16 + uint32(data[3])<<24)
	}

	return i
}

func Byte32Uint32(data []byte, endian int) uint32 {
	var i uint32
	if UBigEndian == endian {
		i = uint32(uint32(data[2]) + uint32(data[1])<<8 + uint32(data[0])<<16)
	}
	
	if ULittleEndian == endian {
		i = uint32(uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16)
	}

	return i
}

func CreateRandomBlock(size uint) []byte {
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