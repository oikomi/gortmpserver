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
	"github.com/oikomi/gortmpserver/server/util"
	"io"
	//"encoding/binary"
	"errors"
)

type ChanBytes chan []byte
//type ChanBytes chan string

type HandShake struct {
	r *bufio.Reader
	w *bufio.Writer
	conn net.Conn
	C0C1 []byte
	C2 []byte
	S0S1S2 []byte
}

//FIXME : ugly code //add error implemet
func (self *HandShake)checkC0C1() error {
	if self.C0C1[0] != 0x03 {
		err := errors.New("C0C1 error")
		return err
	}
	
	return nil
}

func NewHandShake(conn net.Conn) *HandShake {
	return &HandShake {
		r : bufio.NewReader(conn),
		w : bufio.NewWriter(conn),
		conn : conn,	
	}
}

func (self *HandShake)readC0C1() error {
	log.Println("readC0C1")
	if self.C0C1 == nil {
		self.C0C1 = make([]byte, C0Length + C1Length)
		if _, err := io.ReadFull(self.conn, self.C0C1); err != nil {
			return err
		}
	}
	log.Println(self.C0C1)
	
	return nil	
}

func (self *HandShake)readC2() error {
	log.Println("readC2")
	if self.C2 == nil {
		self.C2 = make([]byte, C2Length)
		if _, err := io.ReadFull(self.conn, self.C2); err != nil {
			log.Println(err)
			return err
		}
	}
	log.Println(self.C2)
	
	return nil	
}


func (self *HandShake)writeS0S1S2() error {
	log.Println("writeS0S1S2")
	self.S0S1S2 = util.GenerateRandomBytes(S0Length + S1Length + S2Length)
	
	self.S0S1S2[0] = 0x03
	
	//binary.BigEndian.PutUint32(s1, uint32(0))
	
	//for i := 0; i < 4; i++ {
	//	s1[4+i] = 0x00
	//}
	log.Println(self.S0S1S2)
	log.Println(len(self.S0S1S2))
	
	if _, err := self.conn.Write(self.S0S1S2); err != nil {
		return err
	}
	
	/*
	_, err := self.w.Write(s1)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	
	err = self.w.Flush()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	*/
	
	
	return nil	
}

func (self *HandShake)DoHandshake() error {
	var err error
	if err = self.readC0C1(); err != nil {
		return err
	}
	
	// FIXME
	if err = self.checkC0C1(); err != nil  {
		return err
	}
	
	if err = self.writeS0S1S2(); err != nil  {
		return err
	}
	
	if err = self.readC2(); err != nil {
		return err
	}
	
	return nil
}