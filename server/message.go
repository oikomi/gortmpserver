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
	"bytes"
	"github.com/golang/glog"
)

// Message type
const (
	SET_CHUNK_SIZE = uint8(1)

	ABORT_MESSAGE = uint8(2)

	ACKNOWLEDGEMENT = uint8(3)

	USER_CONTROL_MESSAGE = uint8(4)

	WINDOW_ACKNOWLEDGEMENT_SIZE = uint8(5)

	SET_PEER_BANDWIDTH = uint8(6)

	AUDIO_TYPE = uint8(8)

	VIDEO_TYPE = uint8(9)

	AGGREGATE_MESSAGE_TYPE = uint8(22)

	SHARED_OBJECT_AMF0 = uint8(19)
	SHARED_OBJECT_AMF3 = uint8(16)

	DATA_AMF0 = uint8(18)
	DATA_AMF3 = uint8(15)

	COMMAND_AMF0 = uint8(20)
	COMMAND_AMF3 = uint8(17) // Keng-die!!! Just ignore one byte before AMF0.
)

type Message struct {
	ChunkStreamID     uint32
	Timestamp         uint32
	Size              uint32
	Type              uint8
	StreamID          uint32
	Buf               *bytes.Buffer
	IsInbound         bool
	AbsoluteTimestamp uint32
}

func NewMessage(csi uint32, t uint8, sid uint32, ts uint32, data []byte) *Message {
	message := &Message{
		ChunkStreamID     : csi,
		Type              : t,
		StreamID          : sid,
		Timestamp         : ts,
		AbsoluteTimestamp : ts,
		Buf               : new(bytes.Buffer),
	}
	if data != nil {
		message.Buf.Write(data)
		message.Size = uint32(len(data))
	}
	return message
}

// The length of remain data to read
func (self *Message) Remain() uint32 {
	if self.Buf == nil {
		return self.Size
	}
	return self.Size - uint32(self.Buf.Len())
}


func (self *Message) Dump() {
	glog.Infof("Message{CID: %d, Type: %d, Timestamp: %d, Size: %d, StreamID: %d, IsInbound: %t, AbsoluteTimestamp: %d}\n",
		self.ChunkStreamID, self.Type, self.Timestamp, self.Size, self.StreamID, self.IsInbound, self.AbsoluteTimestamp)
}