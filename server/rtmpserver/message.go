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
	MessageType byte
	PayloadLength uint32
	TimestampDelta uint32
	MessageStreamId uint32
	Timestamp uint64
}

type Message struct {
	Header *MessageHeader
	Payload []byte
}

func NewMessage() *Message {
	return &Message {
		
	}
}

