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
	"net"
	"github.com/golang/glog"
)

type RtmpServer struct {
	cfg     *RtmpServerConfig
	server  *net.TCPListener
}

func NewRtmpServer(cfg *RtmpServerConfig, server *net.TCPListener) *RtmpServer {
	return &RtmpServer {
		cfg    : cfg,
		server : server,
	}
}

func (self *RtmpServer)ServerLoop() {
	glog.Info("ServerLoop")
	for {
		conn, err := self.server.Accept()
		if err != nil {
			glog.Error(err.Error())
			return
		}
		session := NewSession(conn)
		go self.handleSession(session)
	}
}

func (self *RtmpServer)handleSession(session *Session) error {
	glog.Info("handleSession")
	
	err := self.doHandShake(session)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	go session.sendLoop()
	go session.readLoop()
	
	return nil
}  

func (self *RtmpServer)doHandShake(session *Session) error {
	glog.Info("doHandShake")
	var err error
	hs := NewHandShake(self, session)
	err = hs.recvC0C1()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	err = hs.sendS0S1S2()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	err = hs.recvC2()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	glog.Info("doHandShake done")
	
	return nil
}

func (self *RtmpServer)parseProtocol(cs *ChunkStream) error {
	

	return nil
}
