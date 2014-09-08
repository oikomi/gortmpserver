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
	"reflect"
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
	
	recvMessage chan *Message
	sendMessage chan *Message
}

func NewRtmpServer(cfg *config.Config) *RtmpServer {
	return &RtmpServer {
		port : cfg.Listen,
		clients : make(ClientTable),
		recvMessage : make(chan *Message),
		sendMessage : make(chan *Message),
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
	
	req := NewRequest()
	
	err = self.ConnectApp(req)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	ack_size := uint32(2.5 * 1000 * 1000)
	err = self.SetWindowAckSize(ack_size, rtmpClient)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
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
		//log.Printf("Accept")
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}

		go self.handleClient(conn)
		
	}
	
	return nil
}

func (self *RtmpServer) ExpectPacket(v interface {}) (msg *Message, err error) {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr {
		err = errors.New("param must be ptr for expect message")
		return
	}
	if rv.IsNil() {
		err = errors.New("param should never be nil")
		return
	}
	if !rv.Elem().CanSet() {
		err = errors.New("param should be settable")
		return
	}

	for {
		if msg, err = self.RecvMessage(); err != nil {
			log.Fatalln(err.Error())
			return
		}
		
		if msg == nil {
			log.Fatalln("msg == nil")
			continue			
		}
		
		var pkt interface {}
		if pkt, err = self.DecodeMessage(msg); err != nil {
			log.Fatalln(err.Error())
			return
		}
		if pkt == nil {
			log.Fatalln("pkt == nil")
			continue
		}

		// check the convertible and convert to the value or ptr value.
		// for example, the v like the c++ code: Msg**v
		pkt_rt := reflect.TypeOf(pkt)
		if pkt_rt.ConvertibleTo(rt) {
			// directly match, the pkt is like c++: Msg**pkt
			// set the v by: *v = *pkt
			rv.Elem().Set(reflect.ValueOf(pkt).Elem())
			return
		}
		if pkt_rt.ConvertibleTo(rt.Elem()) {
			// ptr match, the pkt is like c++: Msg*pkt
			// set the v by: *v = pkt
			rv.Elem().Set(reflect.ValueOf(pkt))
			return
		}
	}

	return
}

func (self *RtmpServer) ConnectApp(req *Request) error {
	log.Println("ConnectApp")
	var err error
	var pkt *ConnectAppPacket

	_, err = self.ExpectPacket(&pkt)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	v, ok := pkt.CommandObject["tcUrl"]
	if !ok {
		err = errors.New("invalid request, must specifies the tcUrl.")
		return err
	}
	req.TcUrl = v.(string)
	if v, ok := pkt.CommandObject["pageUrl"]; ok {
		req.PageUrl = v.(string)
	}
	if v, ok := pkt.CommandObject["swfUrl"]; ok {
		req.SwfUrl = v.(string)
	}
	if v, ok := pkt.CommandObject["objectEncoding"]; ok {
		req.ObjectEncoding = v.(float64)
	}
	
	log.Println(req.TcUrl)
	log.Println(req.PageUrl)

	err = req.DiscoveryApp()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	return nil	
	
}

func (self *RtmpServer) SetWindowAckSize(ackSize uint32, rtmpClient *RtmpClient)  error {
	pkt := NewSetWindowAckSizePacket()
	n, buf, err := pkt.Encode()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	log.Println("---------")
	log.Println(buf.Len())
	log.Println(n)
	
	return nil	
	
	
	
}

/*
func (self *RtmpServer) EncodeMessage(pkt *IPkgEncode) (*Message, error) {
	msg := NewMessage()

	cid = pkt.GetCid()

	size := pkt.GetSize()

	b := make([]byte, size)
	buf := NewAmfBuffer(b)
	if err = pkt.Encode(buf.Buf); err != nil {
		return
	}

	msg.Header.MessageType = pkt.GetMessageType()
	msg.Header.PayloadLength = uint32(size)
	msg.Payload = b

	return
}
*/

func (self *RtmpServer) RecvMessage() (*Message, error) {
	if msg, ok := <- self.recvMessage; ok {
		return msg, nil
	}

	err := errors.New("recv msg error")
	return nil, err
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
	
	if err = self.readChunkMsgHeader(rtmpClient, chunkStream); err != nil {
		log.Fatalln(err.Error())
		return err
	}


	if err = self.readChunkMsgData(rtmpClient, chunkStream); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	//log.Println(chunkStream.Msg.Payload)	
	
	if err = self.onRecvMessage(chunkStream.Msg); err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	self.recvMessage <- chunkStream.Msg
		
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

func (self *RtmpServer) DecodeMessage(msg *Message) (interface {}, error) {
	if msg == nil || msg.Payload == nil {
		return nil, errors.New("DecodeMessage failed")
	}

	pkt, err := self.DecodePacket(msg.Header, msg.Payload)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	
	return pkt, err
}

func (self *RtmpServer) DecodePacket(header *MessageHeader, payload []byte) (interface {}, error) {
	var pkt IPkgDecode = nil
	var err error
	var cmd string
	buf := amf.NewAmfBuffer(payload)
	if header.IsAmf0Command() {
		cmd, err = amf.ReadString(buf.Buf)
		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		}
			
		switch cmd {
		case AMF0_COMMAND_CONNECT:
			pkt = NewConnectAppPacket()
			err = pkt.FillPacket(buf.Buf)
			if err != nil {
				log.Fatalln(err.Error())
				return nil, err
			}
			
			err = pkt.Verify()	
			if err != nil {
				log.Fatalln(err.Error())
				return nil, err
			}		

		
		}
		
	}

	return pkt, nil
}




