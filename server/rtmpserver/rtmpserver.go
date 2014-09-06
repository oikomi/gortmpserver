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
	"net"
	"log"
	"errors"
	"math"
	"github.com/oikomi/gortmpserver/server/config"
	"github.com/oikomi/gortmpserver/server/handshake"
	"github.com/oikomi/gortmpserver/server/util"
	"github.com/oikomi/gortmpserver/server/amf"
)

type ClientTable map[net.Conn]*RtmpClient

type RtmpServer struct {
	port string
	listener net.Listener
	clients  ClientTable
}

func NewRtmpServer(cfg *config.Config) *RtmpServer {
	return &RtmpServer {
		port : cfg.Listen,
		clients : make(ClientTable),
	}
	
}

func (self *RtmpServer)handleClient(conn net.Conn) error {
	rtmpClient := NewRtmpClient(conn)
	self.clients[conn] = rtmpClient
	hs := handshake.NewHandShake(conn)
	err := hs.DoHandshake()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	self.startMessagePump(rtmpClient)
	
	return nil
	
}

func (self *RtmpServer)Listen() error {
	var err error
	self.listener, err = net.Listen("tcp", self.port)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	for {
		conn, err := self.listener.Accept()
		log.Printf("Accept")
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}

		go self.handleClient(conn)
		
	}
	
	return nil
}

func (self *RtmpServer)startMessagePump(rtmpClient *RtmpClient) {
	go self.recvMsg(rtmpClient)
	go self.sendMsg(rtmpClient)
}

func (self *RtmpServer)recvMsg(rtmpClient *RtmpClient) error {
	chunkStream := NewChunkStream()
	err := self.readChunkBasicHeader(rtmpClient, chunkStream)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	//log.Println(chunkStream.BasicHeader.Fmt)
	//log.Println(chunkStream.BasicHeader.Cid)
	//log.Println(chunkStream.BasicHeader.Size)
	
	if err = self.readChunkMsgHeader(rtmpClient, chunkStream); err != nil {
		log.Fatalln(err.Error())
		return err
	}


	if err = self.readChunkMsgData(rtmpClient, chunkStream); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	log.Println(chunkStream.Msg.Payload)	
	
	if err = self.onRecvMessage(chunkStream.Msg); err != nil {
		log.Fatalln(err.Error())
		return err
	}
		
	return nil
}

func (self *RtmpServer)sendMsg(rtmpClient *RtmpClient) error {
	
	return nil
}

func (self *RtmpServer) readChunkBasicHeader(rtmpClient *RtmpClient, chunkStream *ChunkStream) error {
	format, err := rtmpClient.ReadByte()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	cid := int(format) & 0x3f
	format = (format >> 6) & 0x03
	chunkStream.BasicHeader.Cid = cid
	chunkStream.BasicHeader.Fmt = format
	chunkStream.BasicHeader.Size = 1

	if cid == 0 {
		cid = 64
		tmp, err := rtmpClient.ReadByte()
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}
		cid += int(tmp)
		chunkStream.BasicHeader.Cid = cid
		chunkStream.BasicHeader.Size = 2	
	} else if cid == 1 {
		cid = 64
		tmp, err := rtmpClient.ReadByte()
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}
		cid += int(tmp)
		tmp, err = rtmpClient.ReadByte()
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}
		cid += int(tmp) * 256
		chunkStream.BasicHeader.Cid = cid
		chunkStream.BasicHeader.Size = 3
	}

	return nil
}

func (self *RtmpServer) readChunkMsgHeader(rtmpClient *RtmpClient, chunkStream *ChunkStream) error {
	fmt := int(chunkStream.BasicHeader.Fmt)
	sizesList := []int{11, 7, 3, 0}
	size := sizesList[fmt]
	
	bs, err := rtmpClient.ReadBytes(size)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	if fmt <= RTMP_FMT_TYPE2 {
		chunkStream.MsgHeader.TimestampDelta = util.ReadUint24(bs[0:3])	
		if chunkStream.ExtendedTimestampFlag = false; chunkStream.MsgHeader.TimestampDelta >= RTMP_EXTENDED_TIMESTAMP {
			chunkStream.ExtendedTimestampFlag = true
		}
		
		if chunkStream.ExtendedTimestampFlag {
			chunkStream.MsgHeader.Timestamp = RTMP_EXTENDED_TIMESTAMP
		} else {
			if fmt == RTMP_FMT_TYPE0 {
				chunkStream.MsgHeader.Timestamp = chunkStream.MsgHeader.TimestampDelta
			} else {
				chunkStream.MsgHeader.Timestamp += chunkStream.MsgHeader.TimestampDelta
			}
		}
		
		if fmt <= RTMP_FMT_TYPE1 {
			chunkStream.MsgHeader.MessageLength = util.ReadUint24(bs[3:6])
			chunkStream.MsgHeader.MessageTypeId = util.ReadUint8(bs[6:7])
		}
		
		if fmt <= RTMP_FMT_TYPE0 {
			chunkStream.MsgHeader.MessageStreamId = util.ReadUint32(bs[7:11])	
		}
	}
	
	if chunkStream.ExtendedTimestampFlag {
		bs, err := rtmpClient.ReadBytes(4)
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}
		chunkStream.ExtendedTimestamp = util.ReadUint32(bs)
		log.Println(chunkStream.ExtendedTimestamp)
	}
	
	if chunkStream.MsgHeader.MessageLength  < 0 {
		err = errors.New("chunkStream.MsgHeader.MessageLength  < 0")
		return err
	}
	
	chunkStream.Msg.Header.MessageType = chunkStream.MsgHeader.MessageTypeId
	chunkStream.Msg.Header.MessageStreamId = chunkStream.MsgHeader.MessageStreamId
	chunkStream.Msg.Header.PayloadLength = chunkStream.MsgHeader.MessageLength
	
	return nil
	
}

func (self *RtmpServer) readChunkMsgData(rtmpClient *RtmpClient, chunkStream *ChunkStream) error {
	var err error
	if chunkStream.ChunkData == nil {
		chunkStream.ChunkData = make([]byte, chunkStream.MsgHeader.MessageLength)
	}
	
	if chunkStream.MsgHeader.MessageLength <= RTMP_DEFAULT_CHUNK_SIZE {
		chunkStream.ChunkData, err = rtmpClient.ReadBytes(int(chunkStream.MsgHeader.MessageLength))
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}
		
		log.Println(chunkStream.ChunkData)
		log.Println(len(chunkStream.ChunkData))	
		
		if chunkStream.Msg.Payload == nil {
			chunkStream.Msg.Payload = make([]byte, chunkStream.MsgHeader.MessageLength)
		}
		
		copy(chunkStream.Msg.Payload , chunkStream.ChunkData)
		chunkStream.MsgCount ++
		
	} else {
			if chunkStream.Msg.Payload == nil {
				//log.Println("chunkStream.Msg.Payload == nil")
				chunkStream.Msg.Payload = make([]byte, chunkStream.MsgHeader.MessageLength)
			}
			num := chunkStream.MsgHeader.MessageLength / RTMP_DEFAULT_CHUNK_SIZE
			lastLength := chunkStream.MsgHeader.MessageLength % RTMP_DEFAULT_CHUNK_SIZE
			
			if lastLength != 0 {
				num ++
			}
			chunkStream.MsgCount = num 
			var i uint32
			tmpChunkStream := NewChunkStream()
			tmpSizesList := []int{11, 7, 3, 0}
			
			for i = 0; i < num; i++ {
				size := chunkStream.MsgHeader.MessageLength - chunkStream.Msg.ReceivedPayloadLength
				size = uint32(math.Min(float64(size), float64(RTMP_DEFAULT_CHUNK_SIZE)))

				tmp, err := rtmpClient.ReadBytes(int(size))
				if err != nil {
					log.Fatalln(err.Error())
					return err
				}
				
				copy(chunkStream.Msg.Payload[chunkStream.Msg.ReceivedPayloadLength:chunkStream.Msg.ReceivedPayloadLength+size], tmp)

				chunkStream.Msg.ReceivedPayloadLength += size
				
				if chunkStream.Msg.ReceivedPayloadLength == chunkStream.MsgHeader.MessageLength {
					//log.Println("chunkStream.Msg.ReceivedPayloadLength == chunkStream.MsgHeader.MessageLength")
					break
				}
				
				err = self.readChunkBasicHeader(rtmpClient, tmpChunkStream)
				if err != nil {
					log.Fatalln(err.Error())
					return err
				}
				
				tmpFmt := int(tmpChunkStream.BasicHeader.Fmt)
				
				tmpSize := tmpSizesList[tmpFmt]
				
				if int(tmpSize) != 0 {	
					discardData, err := rtmpClient.ReadBytes(int(tmpSize))
					if err != nil {
						log.Fatalln(err.Error())
						log.Fatalln(discardData)
						return err
					}
				}			
								
			}
	}
	
	return nil
}

func (self *RtmpServer) onRecvMessage(msg *Message) error {
	self.DecodeMessage(msg)
	return nil
}

func (self *RtmpServer)DecodeMessage(msg *Message) {
	self.DecodePacket(msg.Header, msg.Payload)
}

func (self *RtmpServer)DecodePacket(header *MessageHeader, payload []byte) (interface {}, error) {
	if header.IsAmf0Command() {
		var cmd string
		var err error
		amf0Codec := amf.NewAmf0Codec(payload)
		if cmd, err = amf0Codec.ReadString(); err != nil {
			return nil, err
		}
		
		log.Println(cmd)
		
		switch cmd {
		case AMF0_COMMAND_CONNECT:
			NewConnectAppPacket()
		/*
		case AMF0_COMMAND_CREATE_STREAM:
			pkt = NewCreateStreamPacket()
		case AMF0_COMMAND_PLAY:
			pkt = NewPlayPacket()
		case AMF0_COMMAND_PUBLISH:
			pkt = NewPublishPacket()
		case AMF0_COMMAND_CLOSE_STREAM:
			pkt = NewCloseStreamPacket()
		case AMF0_COMMAND_RELEASE_STREAM:
			pkt = NewFMLEStartPacket()
		case AMF0_COMMAND_FC_PUBLISH:
			pkt = NewFMLEStartPacket()
		case AMF0_COMMAND_UNPUBLISH:
			pkt = NewFMLEStartPacket()
		*/
		}
		
	} 
	/*
	else if header.IsWindowAcknowledgementSize() {
		pkt =NewSetWindowAckSizePacket()
	} else if header.IsUserControlMessage() {
		pkt = NewUserControlPacket()
	} else if header.IsSetChunkSize() {
		pkt = NewSetChunkSizePacket()
	}
	*/



	return nil, nil
}



