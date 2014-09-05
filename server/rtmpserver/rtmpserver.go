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
	"github.com/oikomi/gortmpserver/server/config"
	"github.com/oikomi/gortmpserver/server/handshake"
	"github.com/oikomi/gortmpserver/server/util"
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
	
	log.Println(chunkStream.BasicHeader.Fmt)
	log.Println(chunkStream.BasicHeader.Cid)
	log.Println(chunkStream.BasicHeader.Size)
	
	err = self.readChunkMsgHeader(rtmpClient, chunkStream)
	if err != nil {
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
		//log.Println(bs[0:3])
		chunkStream.MsgHeader.TimestampDelta = util.ReadUint24(bs[0:3])	
		if chunkStream.ExtendedTimestampFlag = false; chunkStream.MsgHeader.Timestamp >= RTMP_EXTENDED_TIMESTAMP {
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
	

	
	return nil
	
}

