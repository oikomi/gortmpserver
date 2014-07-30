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

package handshake

import (
	"net"
	"log"
	"bufio"
	//"bytes"
	//"encoding/binary"
)

type ChanBytes chan []byte
//type ChanBytes chan string

type HandShake struct {
	r *bufio.Reader
	w *bufio.Writer
	conn net.Conn
	C0 ChanBytes
	C1 ChanBytes
	S0 ChanBytes
	S1 ChanBytes
}

func NewHandShake(conn net.Conn) *HandShake {
	return &HandShake {
		r : bufio.NewReader(conn),
		w : bufio.NewWriter(conn),
		conn : conn,
		C0 : make(ChanBytes, C0Length),
		C1 : make(ChanBytes, C1Length),
		S0 : make(ChanBytes, S0Length),
		S1 : make(ChanBytes, S1Length),
	}
}

func (self *HandShake)readC0() error {
	tmp := make([]byte, C0Length)
	n, err := self.r.Read(tmp)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	log.Println(n)

	self.C0 <- tmp
	
	return nil	
}

func (self *HandShake)readC1() error {
	tmp := make([]byte, C1Length)
	n, err := self.r.Read(tmp)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	log.Println(n)

	self.C1 <- tmp
	
	return nil	
}

func (self *HandShake)writeS0() error {
	log.Println("writeS0")
	n, err := self.w.Write([]byte{3})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	log.Println(n)
	
	return nil	
}

func (self *HandShake)handShakeEvent() {
	for {
		select {
		case  <-self.C0: //|| message2 := <-self.C1:
			self.writeS0()

		}
	}
}

func (self *HandShake)DoHandshake() error {
	go self.handShakeEvent()
	
	self.readC0()
	self.readC1()
	
	return nil
}

