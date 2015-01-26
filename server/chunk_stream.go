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

package server

import (
	"io"
	"bufio"
	"encoding/binary"
	"github.com/golang/glog"
	"github.com/oikomi/gortmpserver/util"
)

const (
	HEADER_FMT_FULL                   = 0x00
	HEADER_FMT_SAME_STREAM            = 0x01
	HEADER_FMT_SAME_LENGTH_AND_STREAM = 0x02
	HEADER_FMT_CONTINUATION           = 0x03
)

const (
	MESSAGE_HEADER_TYPE0_LENGTH       = 11
	MESSAGE_HEADER_TYPE1_LENGTH       = 7
	MESSAGE_HEADER_TYPE2_LENGTH       = 3
	MESSAGE_HEADER_TYPE3_LENGTH       = 0
)

type BasicHeader struct {
	fmt            uint8
	chunkStreamID  uint32
}

type MesageHeader struct {
	timestamp        uint32
	messageLength    uint32
	messageTypeID    uint8
	messageStreamID  uint32
}

type ChunkStream struct {
	bh                BasicHeader
	mh                MesageHeader
	extendedTimestamp uint32
	data              []byte
	lastChunk         *ChunkStream
}

func NewChunkStream() *ChunkStream {
	return &ChunkStream {
	}
}

func (self *ChunkStream)readBasicHeader(br *bufio.Reader) error {
	var b byte
	b, err := br.ReadByte()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	n := 1
	self.bh.fmt = uint8(b >> 6)
	b = b & 0x3f
	switch b {
	case 0:
		// Chunk stream IDs 64-319 can be encoded in the 2-byte version of this
		// field. ID is computed as (the second byte + 64).
		b, err = br.ReadByte()
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		n += 1
		self.bh.chunkStreamID = uint32(64) + uint32(b)
	case 1:
		// Chunk stream IDs 64-65599 can be encoded in the 3-byte version of
		// this field. ID is computed as ((the third byte)*256 + the second byte
		// + 64).
		b, err = br.ReadByte()
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		n += 1
		self.bh.chunkStreamID = uint32(64) + uint32(b)
		b, err = br.ReadByte()
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		n += 1
		self.bh.chunkStreamID += uint32(b) * 256
	default:
		// Chunk stream IDs 2-63 can be encoded in the 1-byte version of this
		// field.
		self.bh.chunkStreamID = uint32(b)
	}

	return nil
}

func (self *ChunkStream) parseMesageHeader(buf []byte) {
	switch self.bh.fmt {
	case HEADER_FMT_FULL:
		glog.Info(buf)
		self.mh.timestamp = util.Byte32Uint32(buf[0:3], util.UBigEndian)
		self.mh.messageLength = util.Byte32Uint32(buf[3:6], util.UBigEndian)
		self.mh.messageTypeID = uint8(buf[7])
		// messageStreamID is littleEndian 
		self.mh.messageStreamID = util.Byte32Uint32(buf[8:11], util.ULittleEndian)
	case HEADER_FMT_SAME_STREAM:
		self.mh.timestamp = util.Byte32Uint32(buf[0:3], util.UBigEndian)
		self.mh.messageLength = util.Byte32Uint32(buf[3:6], util.UBigEndian)
		self.mh.messageTypeID = uint8(buf[7])
	case HEADER_FMT_SAME_LENGTH_AND_STREAM:
		self.mh.timestamp = util.Byte32Uint32(buf[0:3], util.UBigEndian)

	case HEADER_FMT_CONTINUATION:

	}
}

func (self *ChunkStream) readMesageHeader(br *bufio.Reader) error {
	switch self.bh.fmt {
	case HEADER_FMT_FULL:
		tmpBuf := make([]byte, MESSAGE_HEADER_TYPE0_LENGTH)
		if _, err := io.ReadFull(br, tmpBuf); err != nil {
			return err
		}
		glog.Info(tmpBuf)
		self.parseMesageHeader(tmpBuf)
	case HEADER_FMT_SAME_STREAM:
		tmpBuf := make([]byte, MESSAGE_HEADER_TYPE1_LENGTH)
		if _, err := io.ReadFull(br, tmpBuf); err != nil {
			return err
		}
		self.parseMesageHeader(tmpBuf)

	case HEADER_FMT_SAME_LENGTH_AND_STREAM:
		tmpBuf := make([]byte, MESSAGE_HEADER_TYPE2_LENGTH)
		if _, err := io.ReadFull(br, tmpBuf); err != nil {
			return err
		}
		self.parseMesageHeader(tmpBuf)

	case HEADER_FMT_CONTINUATION:

	}
	//fixme  
	
	if (self.bh.fmt != HEADER_FMT_CONTINUATION && self.mh.timestamp >= 0xffffff) {// ||
		//(self.bh.fmt == HEADER_FMT_CONTINUATION && lastheader != nil && lastheader.ExtendedTimestamp > 0) {
		tmpBuf := make([]byte, 4)
		if _, err := io.ReadFull(br, tmpBuf); err != nil {
			return err
		}
		self.extendedTimestamp = binary.BigEndian.Uint32(tmpBuf)
	} else {
		self.extendedTimestamp = 0
	}
	
	return nil
}

func (self *ChunkStream)ReadChunkStream(br *bufio.Reader) error {
	err := self.readBasicHeader(br)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	err = self.readMesageHeader(br)
	if err != nil {
		glog.Error(err.Error())
		return err
	}

	return nil
}
