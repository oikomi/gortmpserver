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

func (self *RtmpServer)handleClient(conn net.Conn) {
	rtmpclient := NewRtmpClient()
	self.clients[conn] = rtmpclient
	hs := handshake.NewHandShake(conn)
	hs.DoHandshake()

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