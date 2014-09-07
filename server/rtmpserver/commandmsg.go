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
	"errors"
	"bytes"
	"log"
	"github.com/oikomi/gortmpserver/server/amf"
)

type ConnectAppPacket struct {
	CommandName string
	TransactionId float64
	CommandObject amf.Object
}

func NewConnectAppPacket() (*ConnectAppPacket) {
	return &ConnectAppPacket {
		CommandName : AMF0_COMMAND_CONNECT,
		//TransactionId : float64(1.0),
		//CommandObject : make(amf.Object),
	}
}

func (self *ConnectAppPacket)FillPacket(buf *bytes.Buffer) error {
	var err error
	self.TransactionId, err = amf.ReadDouble(buf)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	self.CommandObject, err = amf.ReadObject(buf)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	return nil
}

func (self *ConnectAppPacket)Verify() error {
	if self.CommandName != AMF0_COMMAND_CONNECT {
		return errors.New("ConnectAppPacket : decode name failed")
	}
	
	if self.TransactionId != 1.0 {
		return errors.New("ConnectAppPacket : decode connect TransactionId failed.")
	}

	if self.CommandObject == nil {
		return errors.New("ConnectAppPacket : decode connect CommandObject failed.")
	}

	return nil
}