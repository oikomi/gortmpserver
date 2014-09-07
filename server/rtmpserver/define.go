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
	"bytes"
)

const RTMP_FMT_TYPE0 = 0
const RTMP_FMT_TYPE1 = 1
const RTMP_FMT_TYPE2 = 2
const RTMP_FMT_TYPE3 = 3

const RTMP_EXTENDED_TIMESTAMP  = 0xFFFFFF

const RTMP_DEFAULT_CHUNK_SIZE = 128
const RTMP_MIN_CHUNK_SIZE = 128
const RTMP_MAX_CHUNK_SIZE = 65536

const RTMP_MSG_SetChunkSize  = 0x01
const RTMP_MSG_AbortMessage  = 0x02
const RTMP_MSG_Acknowledgement  = 0x03
const RTMP_MSG_UserControlMessage  = 0x04
const RTMP_MSG_WindowAcknowledgementSize  = 0x05
const RTMP_MSG_SetPeerBandwidth  = 0x06
const RTMP_MSG_EdgeAndOriginServerCommand  = 0x07
const RTMP_MSG_AMF3CommandMessage = 17 // = 0x11
const RTMP_MSG_AMF0CommandMessage = 20 // = 0x14
const RTMP_MSG_AMF0DataMessage = 18 // = 0x12
const RTMP_MSG_AMF3DataMessage = 15 // = 0x0F
const RTMP_MSG_AMF3SharedObject = 16 // = 0x10
const RTMP_MSG_AMF0SharedObject = 19 // = 0x13
const RTMP_MSG_AudioMessage = 8 // = 0x08
const RTMP_MSG_VideoMessage = 9 // = 0x09
const RTMP_MSG_AggregateMessage = 22 // = 0x16


const AMF0_COMMAND_CONNECT = "connect"
const AMF0_COMMAND_CREATE_STREAM = "createStream"
const AMF0_COMMAND_CLOSE_STREAM = "closeStream"
const AMF0_COMMAND_PLAY = "play"
const AMF0_COMMAND_PAUSE = "pause"
const AMF0_COMMAND_ON_BW_DONE = "onBWDone"
const AMF0_COMMAND_ON_STATUS = "onStatus"
const AMF0_COMMAND_RESULT = "_result"
const AMF0_COMMAND_ERROR = "_error"
const AMF0_COMMAND_RELEASE_STREAM = "releaseStream"
const AMF0_COMMAND_FC_PUBLISH = "FCPublish"
const AMF0_COMMAND_UNPUBLISH = "FCUnpublish"
const AMF0_COMMAND_PUBLISH = "publish"
const AMF0_DATA_SAMPLE_ACCESS = "|RtmpSampleAccess"
const AMF0_DATA_SET_DATAFRAME = "@setDataFrame"
const AMF0_DATA_ON_METADATA = "onMetaData"


type IPkg interface {
	Verify() error
	FillPacket(buf *bytes.Buffer) error
}