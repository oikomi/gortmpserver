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

package rtmpserver

import (
	"encoding/binary"
	"errors"
	"bytes"
	"log"
	"os"
	//"github.com/oikomi/gortmpserver/server/amf"
)


type SetWindowAckSizePacket struct {
	SetWindowAckSizeChunkStream *ChunkStream
	AcknowledgementWindowSize uint32
}
func NewSetWindowAckSizePacket() (*SetWindowAckSizePacket) {
	return &SetWindowAckSizePacket{
		SetWindowAckSizeChunkStream : NewChunkStream(),
	}
}

/*
func (self *SetWindowAckSizePacket) GetCid() int {
	return RTMP_CID_ProtocolControl
}
func (self *SetWindowAckSizePacket) GetMessageType() byte {
	return RTMP_MSG_WindowAcknowledgementSize
}
func (self *SetWindowAckSizePacket) GetSize() int {
	return 4
}
*/

// Encode header into io.Writer
func (self *SetWindowAckSizePacket) Encode() (n int, wbuf *bytes.Buffer, err error) {
	self.SetWindowAckSizeChunkStream.BasicHeader.Fmt = 0
	self.SetWindowAckSizeChunkStream.BasicHeader.Cid = RTMP_CID_ProtocolControl
	self.SetWindowAckSizeChunkStream.BasicHeader.Size = 1
	self.SetWindowAckSizeChunkStream.MsgHeader.Timestamp = 0
	self.SetWindowAckSizeChunkStream.MsgHeader.MessageTypeId = RTMP_MSG_WindowAcknowledgementSize
	self.SetWindowAckSizeChunkStream.MsgHeader.MessageLength = 4
	self.SetWindowAckSizeChunkStream.MsgHeader.MessageStreamId = 0
	wbuf = bytes.NewBuffer(make([]byte, 16))
	
	switch {
	case self.SetWindowAckSizeChunkStream.BasicHeader.Cid <= 63:
		err = wbuf.WriteByte(byte((self.SetWindowAckSizeChunkStream.BasicHeader.Fmt << 6) | 
			byte(self.SetWindowAckSizeChunkStream.BasicHeader.Cid)))
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 1
	case self.SetWindowAckSizeChunkStream.BasicHeader.Cid <= 319:
		err = wbuf.WriteByte(self.SetWindowAckSizeChunkStream.BasicHeader.Fmt << 6)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 1
		err = wbuf.WriteByte(byte(self.SetWindowAckSizeChunkStream.BasicHeader.Cid - 64))
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 1
	case self.SetWindowAckSizeChunkStream.BasicHeader.Cid <= 65599:
		err = wbuf.WriteByte((self.SetWindowAckSizeChunkStream.BasicHeader.Fmt << 6) | 0x01)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 1
		tmp := uint16(self.SetWindowAckSizeChunkStream.BasicHeader.Cid - 64)
		err = binary.Write(wbuf, binary.BigEndian, &tmp)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 2
	default:
		return n, wbuf, errors.New("Unsupport chunk stream ID large then 65599")
	}
	tmpBuf := make([]byte, 4)
	var m int
	switch self.SetWindowAckSizeChunkStream.BasicHeader.Fmt {
	case HEADER_FMT_FULL:
		// Timestamp
		binary.BigEndian.PutUint32(tmpBuf, self.SetWindowAckSizeChunkStream.MsgHeader.Timestamp)
		m, err = wbuf.Write(tmpBuf[1:])
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += m
		// Message Length
		binary.BigEndian.PutUint32(tmpBuf, self.SetWindowAckSizeChunkStream.MsgHeader.MessageLength)
		m, err = wbuf.Write(tmpBuf[1:])
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += m
		// Message Type
		err = wbuf.WriteByte(self.SetWindowAckSizeChunkStream.MsgHeader.MessageTypeId)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 1
		// Message Stream ID
		err = binary.Write(wbuf, binary.LittleEndian, &(self.SetWindowAckSizeChunkStream.MsgHeader.MessageStreamId))
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 4
	case HEADER_FMT_SAME_STREAM:
		// Timestamp
		binary.BigEndian.PutUint32(tmpBuf, self.SetWindowAckSizeChunkStream.MsgHeader.Timestamp)
		m, err = wbuf.Write(tmpBuf[1:])
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += m
		// Message Length
		binary.BigEndian.PutUint32(tmpBuf, self.SetWindowAckSizeChunkStream.MsgHeader.MessageLength)
		m, err = wbuf.Write(tmpBuf[1:])
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += m
		// Message Type
		err = wbuf.WriteByte(self.SetWindowAckSizeChunkStream.MsgHeader.MessageTypeId)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += 1
	case HEADER_FMT_SAME_LENGTH_AND_STREAM:
		// Timestamp
		binary.BigEndian.PutUint32(tmpBuf, self.SetWindowAckSizeChunkStream.MsgHeader.Timestamp)
		m, err = wbuf.Write(tmpBuf[1:])
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		n += m
	case HEADER_FMT_CONTINUATION:
	}

	// Type 3 chunks MUST NOT have Extended timestamp????
	// Todo: Test with FMS
	// if header.Timestamp >= 0xffffff && header.Fmt != HEADER_FMT_CONTINUATION {
	/*
	if self.SetWindowAckSizeChunkStream.MsgHeader.Timestamp >= 0xffffff {
		// Extended Timestamp
		err = binary.Write(wbuf, binary.BigEndian, &(header.ExtendedTimestamp))
		if err != nil {
			return
		}
		n += 4
	}
	*/
	log.Println(wbuf.Len())
	//log.Printf("%s", wbuf)
	wbuf.WriteTo(os.Stdout)
	
	return
}

