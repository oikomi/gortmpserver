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
	"github.com/golang/glog"
	"github.com/oikomi/gortmpserver/libnet"
)

type RtmpServer struct {
	cfg     *RtmpServerConfig
	server  *libnet.Server
}

func NewRtmpServer(cfg *RtmpServerConfig, server *libnet.Server) *RtmpServer {
	return &RtmpServer {
		cfg    : cfg,
		server : server,
	}
}

func (self *RtmpServer)ServerLoop() {
	self.server.Serve(func(session *libnet.Session) {
		glog.Info("client ", session.Conn().RemoteAddr().String(), " | in")
		
		go self.handleSession(session)
	})
}

func (self *RtmpServer)handleSession(session *libnet.Session) error {
	glog.Info("handleSession")
	var err error
	err = self.doHandShake(session)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	session.Process(func(msg *libnet.InBuffer) error {
		//glog.Info(string(msg.Data))
		
		err := self.parseProtocol(msg.Data, session)
		if err != nil {
			glog.Error(err.Error())
		}
		
		return nil
	})
	
	return nil
}  

func (self *RtmpServer)doHandShake(session *libnet.Session) error {
	glog.Info("doHandShake")
	var err error
	hs := NewHandShake(self, session)
	err = hs.recvC0C1()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	return nil
}

func (self *RtmpServer)parseProtocol(cmd []byte, session *libnet.Session) error {
	return nil
}