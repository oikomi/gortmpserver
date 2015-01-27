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

import(
	"io"
	"net"
	"bufio"
	"bytes"
	"errors"
	"github.com/golang/glog"
)

const (
	MAX_TIMESTAMP                       = uint32(2000000000)
	AUTO_TIMESTAMP                      = uint32(0XFFFFFFFF)
	DEFAULT_HIGH_PRIORITY_BUFFER_SIZE   = 2048
	DEFAULT_MIDDLE_PRIORITY_BUFFER_SIZE = 128
	DEFAULT_LOW_PRIORITY_BUFFER_SIZE    = 64
	DEFAULT_CHUNK_SIZE                  = uint32(128)
	DEFAULT_WINDOW_SIZE                 = 2500000
	DEFAULT_CAPABILITIES                = float64(15)
	DEFAULT_AUDIO_CODECS                = float64(4071)
	DEFAULT_VIDEO_CODECS                = float64(252)
	FMS_CAPBILITIES                     = uint32(255)
	FMS_MODE                            = uint32(2)
	SET_PEER_BANDWIDTH_HARD             = byte(0)
	SET_PEER_BANDWIDTH_SOFT             = byte(1)
	SET_PEER_BANDWIDTH_DYNAMIC          = byte(2)
)

// Errors
var (
	OldChunkStreamNotExistError    = errors.New("OldChunkStreamNotExist")
)

type Session struct {
	conn             net.Conn
	br               *bufio.Reader
	bw               *bufio.Writer
	closed           bool
	inChunkStreams   map[uint32]*ChunkStream
	outChunkStreams  map[uint32]*ChunkStream
	
	// Chunk size
	inChunkSize      uint32
	outChunkSize     uint32
	outChunkSizeTemp uint32
	
	// Window size
	inWindowSize  uint32
	outWindowSize uint32
	// Bytes counter(For window ack)
	inBytes  uint32
	outBytes uint32
}

func NewSession(conn net.Conn) *Session {
	return &Session {
		conn             : conn,
		br               : bufio.NewReader(conn),
		bw               : bufio.NewWriter(conn),
		outChunkStreams  : make(map[uint32]*ChunkStream),
		inChunkStreams   : make(map[uint32]*ChunkStream),
		inChunkSize      : DEFAULT_CHUNK_SIZE,
		outChunkSize     : DEFAULT_CHUNK_SIZE,
		inWindowSize     : DEFAULT_WINDOW_SIZE,
		outWindowSize    : DEFAULT_WINDOW_SIZE,
	}
}


func (self *Session)bufRead(b []byte) error {
	_, err := self.br.Read(b)
	
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	return nil
}

func (self *Session)bufWrite(b []byte) error {
	_, err := self.bw.Write(b)
	
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	err = self.bw.Flush()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	return nil
}

func (self *Session)parseMessage(message *Message) error {
	tmpBuf := make([]byte, 4)
	var err error
	var subType byte
	var dataSize uint32
	var timestamp uint32
	var timestampExt byte
	if message.Type == AGGREGATE_MESSAGE_TYPE {
		//todo 
		
	} else {
		switch message.ChunkStreamID {
		case CS_ID_PROTOCOL_CONTROL:
			switch message.Type {
			case SET_CHUNK_SIZE:
				conn.invokeSetChunkSize(message)
			case ABORT_MESSAGE:
				conn.invokeAbortMessage(message)
			case ACKNOWLEDGEMENT:
				conn.invokeAcknowledgement(message)
			case USER_CONTROL_MESSAGE:
				conn.invokeUserControlMessage(message)
			case WINDOW_ACKNOWLEDGEMENT_SIZE:
				conn.invokeWindowAcknowledgementSize(message)
			case SET_PEER_BANDWIDTH:
				conn.invokeSetPeerBandwidth(message)
			default:
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE,
					"Unkown message type %d in Protocol control chunk stream!\n", message.Type)
			}
		case CS_ID_COMMAND:
			if message.StreamID == 0 {
				cmd := &Command{}
				var err error
				var transactionID float64
				var object interface{}
				switch message.Type {
				case COMMAND_AMF3:
					cmd.IsFlex = true
					_, err = message.Buf.ReadByte()
					if err != nil {
						logger.ModulePrintln(logHandler, log.LOG_LEVEL_WARNING,
							"Read first in flex commad err:", err)
						return
					}
					fallthrough
				case COMMAND_AMF0:
					cmd.Name, err = amf.ReadString(message.Buf)
					if err != nil {
						logger.ModulePrintln(logHandler, log.LOG_LEVEL_WARNING,
							"AMF0 Read name err:", err)
						return
					}
					transactionID, err = amf.ReadDouble(message.Buf)
					if err != nil {
						logger.ModulePrintln(logHandler, log.LOG_LEVEL_WARNING,
							"AMF0 Read transactionID err:", err)
						return
					}
					cmd.TransactionID = uint32(transactionID)
					for message.Buf.Len() > 0 {
						object, err = amf.ReadValue(message.Buf)
						if err != nil {
							logger.ModulePrintln(logHandler, log.LOG_LEVEL_WARNING,
								"AMF0 Read object err:", err)
							return
						}
						cmd.Objects = append(cmd.Objects, object)
					}
				default:
					logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE,
						"Unkown message type %d in Command chunk stream!\n", message.Type)
				}
				conn.invokeCommand(cmd)
			} else {
				conn.handler.OnReceived(conn, message)
			}
		default:
			conn.handler.OnReceived(conn, message)
		}
	}
	return nil

}

func (self *Session)readLoop() error {
	glog.Info("readLoop")
	var err error
	for !self.closed {
		cs := NewChunkStream()
		cs.readChunkStream(self.br)
		oldChunkStream, found := self.inChunkStreams[cs.bh.chunkStreamID]
		if !found || oldChunkStream == nil {
			self.inChunkStreams[cs.bh.chunkStreamID] = cs
		}
		var absoluteTimestamp uint32
		var message *Message
		switch cs.bh.fmt {
		case HEADER_FMT_FULL:
			cs.lastChunk = cs
			absoluteTimestamp = cs.mh.timestamp
		case HEADER_FMT_SAME_STREAM:
			// A new message with same stream ID
			if oldChunkStream == nil {
				// error
				glog.Error(OldChunkStreamNotExistError)
				return OldChunkStreamNotExistError
			} else {
				cs.mh.messageStreamID = oldChunkStream.mh.messageStreamID
			}
			cs.lastChunk = oldChunkStream
			absoluteTimestamp = oldChunkStream.lastInAbsoluteTimestamp + cs.mh.timestamp
		case HEADER_FMT_SAME_LENGTH_AND_STREAM:
			// A new message with same stream ID, message length and message type
			if oldChunkStream == nil {
				glog.Error(OldChunkStreamNotExistError)
				return OldChunkStreamNotExistError
			}
			cs.mh.messageStreamID = oldChunkStream.mh.messageStreamID
			cs.mh.messageLength = oldChunkStream.mh.messageLength
			cs.mh.messageTypeID = oldChunkStream.mh.messageTypeID
			cs.lastChunk = oldChunkStream
			absoluteTimestamp = oldChunkStream.lastInAbsoluteTimestamp + cs.mh.timestamp
		case HEADER_FMT_CONTINUATION:
			if oldChunkStream.receivedMessage != nil {
				// Continuation the previous unfinished message
				message = oldChunkStream.receivedMessage
			}
			if oldChunkStream == nil {
				glog.Error(OldChunkStreamNotExistError)
				return OldChunkStreamNotExistError
			} else {
				cs.mh.messageStreamID = oldChunkStream.mh.messageStreamID
				cs.mh.messageLength = oldChunkStream.mh.messageLength
				cs.mh.messageTypeID = oldChunkStream.mh.messageTypeID
				cs.mh.timestamp = oldChunkStream.mh.timestamp
			}
			cs.lastChunk = oldChunkStream
			absoluteTimestamp = oldChunkStream.lastInAbsoluteTimestamp
		}
		
		cs.lastInAbsoluteTimestamp = absoluteTimestamp
		self.inChunkStreams[cs.bh.chunkStreamID] = cs
		
		cs.dump()
		
		if message == nil {
			// New message
			message = &Message {
				ChunkStreamID:     cs.bh.chunkStreamID,
				Type:              cs.mh.messageTypeID,
				Timestamp:         cs.realTimestamp(),
				Size:              cs.mh.messageLength,
				StreamID:          cs.mh.messageStreamID,
				Buf:               new(bytes.Buffer),
				IsInbound:         true,
				AbsoluteTimestamp: absoluteTimestamp,
			}
		}
		
		remain := message.Remain()
		var n64 int64
		if remain <= self.inChunkSize {
			// One chunk message
			glog.Info("One chunk message")
			for {
				n64, err = io.CopyN(message.Buf, self.br, int64(remain))
				if err != nil {
					glog.Error(err.Error())
					return err
				}
				if err == nil {
					self.inBytes += uint32(n64)
					if remain <= uint32(n64) {
						break
					} else {
						remain -= uint32(n64)
						continue
					}
				}
			}
			// Finished message
			//glog.Info(len(message.Buf.Bytes()))
			//glog.Info(message.Buf.Bytes())
			self.parseMessage(message)
			cs.receivedMessage = nil
		} else {
			// Unfinish
			remain = self.inChunkSize
			for {
				n64, err = io.CopyN(message.Buf, self.br, int64(remain))
				if err != nil {
					glog.Error(err.Error())
					return err
				}

				self.inBytes += uint32(n64)
				if remain <= uint32(n64) {
					break
				} else {
					remain -= uint32(n64)
					continue
				}
				break
			}
			cs.receivedMessage = message
		}

	}
	return nil
}

func (self *Session)sendLoop() error {
	
	return nil
}
