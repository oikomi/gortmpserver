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

type MessageHeader struct {
	MessageType uint8
	PayloadLength uint32
	TimestampDelta uint32
	MessageStreamId uint32
	Timestamp uint64
}

func (self *MessageHeader) IsAmf0Command() (bool) {
	return self.MessageType == RTMP_MSG_AMF0CommandMessage
}
func (self *MessageHeader) IsAmf3Command() (bool) {
	return self.MessageType == RTMP_MSG_AMF3CommandMessage
}
func (self *MessageHeader) IsAmf0Data() (bool) {
	return self.MessageType == RTMP_MSG_AMF0DataMessage
}
func (self *MessageHeader) IsAmf3Data() (bool) {
	return self.MessageType == RTMP_MSG_AMF3DataMessage
}
func (self *MessageHeader) IsWindowAcknowledgementSize() (bool) {
	return self.MessageType == RTMP_MSG_WindowAcknowledgementSize
}
func (self *MessageHeader) IsSetChunkSize() (bool) {
	return self.MessageType == RTMP_MSG_SetChunkSize
}
func (self *MessageHeader) IsUserControlMessage() (bool) {
	return self.MessageType == RTMP_MSG_UserControlMessage
}
func (self *MessageHeader) IsVideo() (bool) {
	return self.MessageType == RTMP_MSG_VideoMessage
}
func (self *MessageHeader) IsAudio() (bool) {
	return self.MessageType == RTMP_MSG_AudioMessage
}
func (self *MessageHeader) IsAggregate() (bool) {
	return self.MessageType == RTMP_MSG_AggregateMessage
}


type Message struct {
	Header *MessageHeader
	ReceivedPayloadLength uint32
	Payload []byte
}

func NewMessage() *Message {
	return &Message {
		Header : new(MessageHeader),
	}
}

