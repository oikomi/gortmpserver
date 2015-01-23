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

type HandShake struct {
	rtmpServer  *RtmpServer
	session     *libnet.Session
}

func NewHandShake(rtmpServer *RtmpServer, session *libnet.Session) *HandShake {
	return &HandShake{
		rtmpServer : rtmpServer,
		session    : session,
	}
}

func (self *HandShake)recvC0C1() error {
	var err error
	err = self.session.ProcessOnce(func(msg *libnet.InBuffer) error {
		glog.Info(string(msg.Data))
		return nil
	})
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	return nil
}
