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
	"bufio"
	"log"
	"errors"
)

type RtmpClient struct {
	conn net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewRtmpClient(conn net.Conn) *RtmpClient {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	return &RtmpClient {
		conn : conn,
		reader : r,
		writer : w,
	}
}

func (self *RtmpClient) ReadByte() (byte, error) {
	b, err := self.reader.ReadByte()
	if err != nil {
		log.Fatalln(err.Error())
		return b, err
	}
	
	return b, nil
}

func (self *RtmpClient) ReadBytes(num int) ([]byte, error) {
	bs := make([]byte, num)
	n, err := self.reader.Read(bs)
	if err != nil {
		log.Fatalln(err.Error())
		return bs, err
	}
	
	if n != num {
		err = errors.New("read wrong num bytes")
		return bs, err
	}
	
	return bs, nil
}
