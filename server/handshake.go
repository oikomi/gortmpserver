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
	"time"
	"bytes"
	"errors"
	"encoding/binary"
	"github.com/golang/glog"
	"github.com/oikomi/gortmpserver/util"
)

const (
	RTMP_HANDSHAKE_C0_LENGTH  =  1
	RTMP_HANDSHAKE_C1_LENGTH  =  1536
	RTMP_HANDSHAKE_C2_LENGTH  =  1536
	RTMP_HANDSHAKE_S0_LENGTH  =  1
	RTMP_HANDSHAKE_S1_LENGTH  =  1536
	RTMP_HANDSHAKE_S2_LENGTH  =  1536
)

const (
	VERSION  = 0x03
)

// Errors
var (
	HandShakeFailedError    = errors.New("HandShake Failed")
)

type HandShake struct {
	rtmpServer  *RtmpServer
	session     *Session
	old         bool
	c0c1        []byte
	s0s1s2      *bytes.Buffer
	c2          []byte
}

func NewHandShake(rtmpServer *RtmpServer, session *Session) *HandShake {
	return &HandShake{
		rtmpServer : rtmpServer,
		session    : session,
		old        : false,
		c0c1       : make([]byte, RTMP_HANDSHAKE_C0_LENGTH + RTMP_HANDSHAKE_C1_LENGTH),
		s0s1s2     : new(bytes.Buffer),
		c2         : make([]byte, RTMP_HANDSHAKE_C2_LENGTH),
	}
}



func (self *HandShake)parseC0C1() error {
	if len(self.c0c1) != RTMP_HANDSHAKE_C0_LENGTH + RTMP_HANDSHAKE_C1_LENGTH {
		glog.Error(HandShakeFailedError)
		return HandShakeFailedError
	}
	//todo : if server cannot identify client version, server should response 0x03 
	if self.c0c1[0] != VERSION {
		glog.Error(HandShakeFailedError)
		return HandShakeFailedError
	}
	if util.Byte42Uint32(self.c0c1[5:9], 0) == 0 {
		glog.Info("old handshake")
		self.old = true
		return nil
	}
	
	//todo : add complex handshake

	return nil
}

func (self *HandShake)recvC0C1() error {
	glog.Info("handshake | recvC0C1")
	err := self.session.bufRead(self.c0c1)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	//glog.Info(self.c0c1)
	
	err = self.parseC0C1()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	return nil
}

func (self *HandShake)genS0S1S2() error {
	if self.old == true {
		//fixme: not sure about these code
		binary.Write(self.s0s1s2, binary.BigEndian, uint8(VERSION))
		epoch1 := time.Now().Unix()
		binary.Write(self.s0s1s2, binary.BigEndian, uint32(epoch1))
		binary.Write(self.s0s1s2, binary.BigEndian, uint32(0))
		b1 := util.CreateRandomBlock(1528)
		binary.Write(self.s0s1s2, binary.BigEndian, b1)
		
		binary.Write(self.s0s1s2, binary.BigEndian, self.c0c1[1:5])
		epoch2 := time.Now().Unix()
		binary.Write(self.s0s1s2, binary.BigEndian, uint32(epoch2))
		binary.Write(self.s0s1s2, binary.BigEndian, self.c0c1[9:1537])
		
		//binary.Write(self.s0s1s2, binary.BigEndian, self.c0c1)
		
		//glog.Info(self.s0s1s2.Bytes())
		//glog.Info(len(self.s0s1s2.Bytes()))
	} else {
		
	}
	
	return nil
}

func (self *HandShake)sendS0S1S2() error {
	glog.Info("handshake | sendS0S1S2")
	err := self.genS0S1S2()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	err = self.session.bufWrite(self.s0s1s2.Bytes())
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	return nil
}

func (self *HandShake)parseC2() error {
	
	return nil
}

func (self *HandShake)recvC2() error {
	glog.Info("handshake | recvC2")
	err := self.session.bufRead(self.c2)
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	//glog.Info(self.c2)
	
	return nil
}
